package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger *slog.Logger

func init() {
	// Initialize structured JSON logging to stdout
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// generateTraceID generates a random 16-byte hex trace ID (W3C TraceContext standard)
func generateTraceID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateSpanID generates a random 8-byte hex span ID
func generateSpanID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// loggingMiddleware captures request context and emits structured JSON access logs
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Trace Context: propagate or generate trace_id
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = generateTraceID()
		}
		spanID := generateSpanID()

		// Add trace headers to response for downstream correlation
		w.Header().Set("X-Trace-Id", traceID)
		w.Header().Set("X-Span-Id", spanID)

		rw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(rw, r)

		duration := time.Since(start).Milliseconds()

		// Emit structured JSON log entry
		logAttrs := []slog.Attr{
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Int64("duration_ms", duration),
			slog.String("client_ip", r.RemoteAddr),
		}

		if rw.statusCode >= 500 {
			logger.LogAttrs(r.Context(), slog.LevelError, "HTTP request failed", logAttrs...)
		} else if rw.statusCode >= 400 {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "HTTP client warning", logAttrs...)
		} else {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "HTTP request completed", logAttrs...)
		}
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	hostname, _ := os.Hostname()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting Day 10 Observability Service",
		slog.String("version", "1.0.0"),
		slog.String("hostname", hostname),
		slog.String("port", port),
	)

	// Root endpoint
	http.HandleFunc("/", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "=== Kubernetes Observability II: Logging & Tracing Demo (Day 10) ===\n"+
			"Hostname : %s\n"+
			"Time     : %s\n\n"+
			"Available Endpoints:\n"+
			"  GET /healthz              -> Health check probe (200 OK)\n"+
			"  GET /api/checkout?user=X  -> Simulates multi-step checkout logs (INFO & WARN)\n"+
			"  GET /api/fail?reason=X    -> Simulates severe backend failures (ERROR)\n"+
			"  GET /api/slow?ms=X        -> Simulates slow queries with high duration_ms\n",
			hostname, time.Now().Format(time.RFC3339))
	}))

	// Health check endpoint
	http.HandleFunc("/healthz", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	}))

	// Simulated multi-step business workflow emitting contextual logs
	http.HandleFunc("/api/checkout", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		traceID := w.Header().Get("X-Trace-Id")
		user := r.URL.Query().Get("user")
		if user == "" {
			user = "anonymous-customer"
		}
		item := r.URL.Query().Get("item")
		if item == "" {
			item = "k8s-handbook"
		}

		// Step 1: Inventory check
		logger.Info("Checking product inventory",
			slog.String("trace_id", traceID),
			slog.String("item", item),
			slog.String("user", user),
			slog.Int("available_stock", 14),
		)

		// Step 2: Payment processing
		time.Sleep(30 * time.Millisecond)
		logger.Info("Payment authorization successful",
			slog.String("trace_id", traceID),
			slog.String("gateway", "stripe-mock"),
			slog.Float64("amount_usd", 49.99),
		)

		// Step 3: Order confirmation
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "confirmed", "item": "%s", "user": "%s", "trace_id": "%s"}`+"\n",
			item, user, traceID)
	}))

	// Simulated failure endpoint generating rich ERROR logs
	http.HandleFunc("/api/fail", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		traceID := w.Header().Get("X-Trace-Id")
		reason := r.URL.Query().Get("reason")
		if reason == "" {
			reason = "database_connection_refused"
		}

		logger.Error("Critical downstream dependency failure",
			slog.String("trace_id", traceID),
			slog.String("error_type", reason),
			slog.String("service", "mysql-primary.prod.svc"),
			slog.String("stack_trace", "main.go:132 -> db/connection.go:48 -> net/dial.go:120"),
			slog.Int("retry_attempts", 3),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Internal Server Error", "reason": "%s", "trace_id": "%s"}`+"\n",
			reason, traceID)
	}))

	// Simulated slow request
	http.HandleFunc("/api/slow", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		traceID := w.Header().Get("X-Trace-Id")
		msStr := r.URL.Query().Get("ms")
		ms := 300
		if val, err := strconv.Atoi(msStr); err == nil && val > 0 {
			ms = val
		}

		time.Sleep(time.Duration(ms) * time.Millisecond)

		logger.Warn("Slow database query detected",
			slog.String("trace_id", traceID),
			slog.Int("query_duration_ms", ms),
			slog.String("table", "audit_events"),
		)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Completed slow request in %d ms (trace: %s)\n", ms, traceID)
	}))

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Error("Server crashed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
