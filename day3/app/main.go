package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	// Get Redis host from environment variable, default to "localhost" (for local dev)
	// In K8s, this will be "redis-service"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "Hello from Go! Running on Pod: %s\n", hostname)
	})

	http.HandleFunc("/incr", func(w http.ResponseWriter, r *http.Request) {
		val, err := rdb.Incr(ctx, "hits").Result()
		if err != nil {
			http.Error(w, "Error connecting to Redis: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Hits: %d\n", val)
	})

	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
