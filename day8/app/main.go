package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// Global slice to hold allocated memory so GC does not collect it
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

	// Root endpoint: displays pod info and current memory stats
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		sinkLock.Lock()
		heldChunks := len(memorySink)
		sinkLock.Unlock()

		resp := fmt.Sprintf(
			"=== Kubernetes Resource & QoS Demo ===\n"+
				"Pod Hostname      : %s\n"+
				"Go Version        : %s\n"+
				"Num CPU (vCPU)    : %d\n"+
				"Num Goroutines    : %d\n"+
				"Allocated Memory  : %d MB (%d bytes)\n"+
				"Heap In-Use       : %d MB (%d bytes)\n"+
				"System Memory     : %d MB (%d bytes)\n"+
				"GC Cycles         : %d\n"+
				"Held 1MB Chunks   : %d\n\n"+
				"--- Available Endpoints ---\n"+
				"GET /healthz              -> Health probe\n"+
				"GET /allocate?mb=<N>      -> Allocate N megabytes (triggers OOMKilled if > limit)\n"+
				"GET /burn-cpu?ms=<ms>&c=<c> -> Burn CPU on 'c' cores for 'ms' milliseconds\n"+
				"GET /free                 -> Release allocated memory back to OS\n",
			hostname,
			runtime.Version(),
			runtime.NumCPU(),
			runtime.NumGoroutine(),
			bToMb(m.Alloc), m.Alloc,
			bToMb(m.HeapInuse), m.HeapInuse,
			bToMb(m.Sys), m.Sys,
			m.NumGC,
			heldChunks,
		)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	})

	// Health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	// Memory stress endpoint: forces actual RSS allocation
	http.HandleFunc("/allocate", func(w http.ResponseWriter, r *http.Request) {
		mbStr := r.URL.Query().Get("mb")
		mb, err := strconv.Atoi(mbStr)
		if err != nil || mb <= 0 {
			http.Error(w, "Invalid or missing 'mb' query parameter (e.g. /allocate?mb=50)", http.StatusBadRequest)
			return
		}

		log.Printf("[ALLOCATE] Beginning allocation of %d MB...", mb)
		sinkLock.Lock()
		defer sinkLock.Unlock()

		// Allocate in 1 MB chunks and write to every page (4096 bytes) so the OS kernel
		// actively faults and commits real physical memory into the container's cgroup.
		for i := 0; i < mb; i++ {
			chunk := make([]byte, 1024*1024)
			for j := 0; j < len(chunk); j += 4096 {
				chunk[j] = byte(i % 256)
			}
			memorySink = append(memorySink, chunk)
			if (i+1)%10 == 0 || i+1 == mb {
				log.Printf("[ALLOCATE] Successfully allocated %d / %d MB", i+1, mb)
			}
			// Small 20ms pause between chunks so logs can be streamed live before OOM
			time.Sleep(20 * time.Millisecond)
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		msg := fmt.Sprintf("Allocated %d MB. Current process Alloc: %d MB, Total held chunks: %d\n",
			mb, bToMb(m.Alloc), len(memorySink))
		w.Write([]byte(msg))
	})

	// CPU stress endpoint: burns CPU cycles
	http.HandleFunc("/burn-cpu", func(w http.ResponseWriter, r *http.Request) {
		msStr := r.URL.Query().Get("ms")
		coresStr := r.URL.Query().Get("c")

		ms := 3000
		if val, err := strconv.Atoi(msStr); err == nil && val > 0 {
			ms = val
		}
		cores := 1
		if val, err := strconv.Atoi(coresStr); err == nil && val > 0 {
			cores = val
		}

		log.Printf("[CPU-BURN] Burning CPU across %d worker(s) for %d ms...", cores, ms)
		stop := make(chan struct{})
		var wg sync.WaitGroup

		for i := 0; i < cores; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-stop:
						return
					default:
						// Tight loop doing math
						_ = 3.1415926535 * 2.7182818284
					}
				}
			}()
		}

		time.Sleep(time.Duration(ms) * time.Millisecond)
		close(stop)
		wg.Wait()

		log.Printf("[CPU-BURN] Finished CPU burn test.")
		fmt.Fprintf(w, "Burned CPU across %d worker(s) for %d ms.\n", cores, ms)
	})

	// Release memory endpoint
	http.HandleFunc("/free", func(w http.ResponseWriter, r *http.Request) {
		sinkLock.Lock()
		count := len(memorySink)
		memorySink = nil
		sinkLock.Unlock()

		runtime.GC()
		debug.FreeOSMemory()

		log.Printf("[FREE] Released %d chunks and ran GC + FreeOSMemory()", count)
		fmt.Fprintf(w, "Freed %d MB chunks. Forced runtime.GC() and FreeOSMemory().\n", count)
	})

	log.Printf("Starting Resource & QoS Demo Server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
