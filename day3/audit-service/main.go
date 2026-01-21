package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type AuditLog struct {
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var entry AuditLog
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// In a real app, we'd save this to a DB or Splunk.
		// Here, we just print to stdout for `kubectl logs` verification.
		hostname, _ := os.Hostname()
		log.Printf("[AUDIT] Pod: %s | Action: %s | Details: %s | Time: %s",
			hostname, entry.Action, entry.Details, entry.Timestamp)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Logged")
	})

	port := "8080"
	log.Printf("Audit Service starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
