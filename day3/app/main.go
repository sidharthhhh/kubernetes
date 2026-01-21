package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Service Discovery: get the Todo API host (Service name)
	todoHost := os.Getenv("TODO_API_HOST")
	if todoHost == "" {
		todoHost = "localhost" // dev mode
	}
	todoPort := os.Getenv("TODO_API_PORT")
	if todoPort == "" {
		todoPort = "8080"
	}

	todoServiceURL := fmt.Sprintf("http://%s:%s", todoHost, todoPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()

		html := fmt.Sprintf(`
		<html>
		<head><title>K8s Day 3 Frontend</title></head>
		<body>
			<h1>Frontend (Pod: %s)</h1>
			<p>Connects to: <strong>%s</strong></p>
			<hr>
			<h2>Create Todo</h2>
			<form action="/create" method="POST">
				<input type="text" name="title" placeholder="Buy milk">
				<button type="submit">Add</button>
			</form>
			<hr>
			<h2>Todo List</h2>
			<div id="todos">Loading...</div>
			
			<script>
				fetch('/todos')
					.then(response => response.json())
					.then(data => {
						const list = document.getElementById('todos');
						if (data.length === 0) {
							list.innerHTML = "No todos found.";
							return;
						}
						let html = "<ul>";
						data.forEach(item => {
							html += "<li>" + item.title + "</li>";
						});
						html += "</ul>";
						list.innerHTML = html;
					})
					.catch(err => document.getElementById('todos').innerHTML = "Error loading todos: " + err);
			</script>
		</body>
		</html>
		`, hostname, todoServiceURL)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	// Proxy: Get Todos
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(todoServiceURL + "/todos")
		if err != nil {
			http.Error(w, "Backend unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	// Proxy: Create Todo
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		title := r.FormValue("title")
		jsonBody := []byte(fmt.Sprintf(`{"title": "%s", "done": false}`, title))

		resp, err := http.Post(todoServiceURL+"/todos", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			http.Error(w, "Backend error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Redirect back home
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	port := "8080"
	log.Printf("Frontend starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
