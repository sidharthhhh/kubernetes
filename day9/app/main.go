package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus Metrics Definitions
var (
	// 1. Counter: Only goes up (or resets to 0 on restart). Ideal for request counts, error counts.
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by status code, method and path.",
		},
		[]string{"method", "path", "status"},
	)

	// 2. Histogram: Samples observations (usually request durations) into configurable buckets.
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	// 3. Gauge: Can go up and down. Ideal for in-flight requests, memory usage, worker pool size.
	activeRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of currently active in-flight HTTP requests.",
		},
	)
)

func init() {
	// Register metrics with the default Prometheus registry
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(activeRequests)
}

// responseWriterInterceptor captures the HTTP status code
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterInterceptor) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// prometheusMiddleware wraps handlers to automatically record metrics
func prometheusMiddleware(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activeRequests.Inc()
		defer activeRequests.Dec()

		start := time.Now()
		rw := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next(rw, r)

		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rw.statusCode)).Inc()
	}
}

func main() {
	hostname, _ := os.Hostname()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Home endpoint
	http.HandleFunc("/", prometheusMiddleware("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "=== Kubernetes Observability Demo (Day 9) ===\n"+
			"Hostname : %s\n"+
			"Time     : %s\n\n"+
			"Available Endpoints:\n"+
			"  GET /metrics       -> Prometheus Metrics Endpoint (Scraped by Prometheus)\n"+
			"  GET /healthz       -> Health check probe (Returns 200 OK)\n"+
			"  GET /api/work      -> Simulates business logic (random latency & occasional 500s)\n"+
			"  GET /api/error     -> Always returns HTTP 500 Internal Server Error\n",
			hostname, time.Now().Format(time.RFC3339))
	}))

	// Health check endpoint
	http.HandleFunc("/healthz", prometheusMiddleware("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	}))

	// Simulated work endpoint with varying latency and realistic HTTP error rate
	http.HandleFunc("/api/work", prometheusMiddleware("/api/work", func(w http.ResponseWriter, r *http.Request) {
		// Simulate random work duration between 10ms and 250ms
		sleepMs := 10 + rand.Intn(240)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		// 10% chance of returning HTTP 500 Internal Server Error
		if rand.Float32() < 0.10 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Simulated intermittent failure", "latency_ms": %d}`+"\n", sleepMs)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success", "processed_by": "%s", "latency_ms": %d}`+"\n", hostname, sleepMs)
	}))

	// Explicit error endpoint for alerting tests
	http.HandleFunc("/api/error", prometheusMiddleware("/api/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error", "message": "Triggered intentional error"}` + "\n"))
	}))

	// Expose the Prometheus /metrics endpoint (includes default Go runtime & our custom metrics)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Server listening on port %s (Prometheus metrics enabled at /metrics)...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
