package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	// connect to redis
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Printf("Connecting to Redis at %s...", redisAddr)
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Basic health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	// CRUD for Todos
	http.HandleFunc("/todos", handleTodos)

	port := "8080"
	log.Printf("Todo API starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// GET: List all todos
		val, err := rdb.Get(ctx, "todos").Result()
		if err == redis.Nil {
			// Key does not exist
			fmt.Fprint(w, "[]")
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, val)

	} else if r.Method == http.MethodPost {
		// POST: Create a new todo
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Ideally we would ensure ID uniqueness etc, but for this lab simpler is better.
		// We will fetch existing list, append, and save back.
		// Not concurrency safe but fine for a lab.

		var todos []Todo
		val, err := rdb.Get(ctx, "todos").Result()
		if err == nil {
			json.Unmarshal([]byte(val), &todos)
		}

		todos = append(todos, todo)

		data, _ := json.Marshal(todos)
		rdb.Set(ctx, "todos", string(data), 0)

		// Call Audit Service (Fire and forget ish)
		log.Printf("DEBUG: Triggering Audit for Todo: %s", todo.Title)
		go func(t Todo) {
			log.Println("DEBUG: Inside Audit Goroutine")
			auditHost := os.Getenv("AUDIT_SERVICE_HOST")
			if auditHost == "" {
				log.Println("Audit service not configured")
				return
			}
			auditPort := os.Getenv("AUDIT_SERVICE_PORT")
			if auditPort == "" {
				auditPort = "80"
			}

			auditURL := fmt.Sprintf("http://%s:%s/log", auditHost, auditPort)
			// log.Printf("DEBUG: Sending audit to %s", auditURL)

			logEntry := map[string]interface{}{
				"action":    "CREATE_TODO",
				"details":   fmt.Sprintf("Created todo: %s", t.Title),
				"timestamp": time.Now(),
			}
			jsonBody, _ := json.Marshal(logEntry)

			resp, err := http.Post(auditURL, "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				log.Printf("Failed to audit: %v", err)
				return
			}
			defer resp.Body.Close()
			log.Println("Audit log sent successfully")
		}(todo)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
