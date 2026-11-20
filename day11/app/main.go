package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	memorySink [][]byte
	sinkLock   sync.Mutex
)

func main() {
	hostname, _ := os.Hostname()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	// Root endpoint: displays pod hostname and available endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "=== Kubernetes Autoscaling Demo (Day 11) ===\n"+
			"Pod Hostname : %s\n"+
			"Current Time : %s\n\n"+
			"Available Endpoints:\n"+
			"  GET /healthz            -> Health check (Liveness / Readiness probe)\n"+
			"  GET /cpu-load?loops=X   -> Generates CPU-bound hashing workload (Default: 300,000 loops)\n"+
			"  GET /memory?mb=X        -> Allocates X megabytes of memory (For VPA / Memory HPA)\n"+
			"  GET /reset              -> Frees allocated memory\n",
			hostname, time.Now().Format(time.RFC3339))
	})

	// Health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	// CPU-intensive workload to trigger HPA scaling
	http.HandleFunc("/cpu-load", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		loopsStr := r.URL.Query().Get("loops")
		loops := 300000 // Default loop count: produces ~50-100ms CPU burn per request
		if val, err := strconv.Atoi(loopsStr); err == nil && val > 0 {
			loops = val
		}

		// Perform cryptographic hashing in a tight loop to burn CPU
		data := []byte("kubernetes-autoscaling-hpa-load-test")
		hash := sha256.Sum256(data)
		for i := 0; i < loops; i++ {
			hash = sha256.Sum256(hash[:])
		}

		duration := time.Since(start)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "CPU load completed on %s: %d hashing rounds in %v\n", hostname, loops, duration)
	})

	// Memory allocation endpoint for memory-based autoscaling or VPA
	http.HandleFunc("/memory", func(w http.ResponseWriter, r *http.Request) {
		mbStr := r.URL.Query().Get("mb")
		mb := 10
		if val, err := strconv.Atoi(mbStr); err == nil && val > 0 {
			mb = val
		}

		sinkLock.Lock()
		for i := 0; i < mb; i++ {
			chunk := make([]byte, 1024*1024)
			for j := 0; j < len(chunk); j += 4096 {
				chunk[j] = byte(i % 256)
			}
			memorySink = append(memorySink, chunk)
		}
		total := len(memorySink)
		sinkLock.Unlock()

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Allocated %d MB on %s. Total held: %d MB\n", mb, hostname, total)
	})

	// Reset memory
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		sinkLock.Lock()
		count := len(memorySink)
		memorySink = nil
		sinkLock.Unlock()

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Freed %d MB on %s\n", count, hostname)
	})

	log.Printf("Autoscaling Demo server running on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
