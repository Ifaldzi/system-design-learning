package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	base := flag.String("base", "http://localhost:8080", "Base URL of the running URL shortener")
	total := flag.Int("n", 10000, "Total requests for the throughput/concurrency stress test")
	workers := flag.Int("c", 50, "Number of concurrent worker goroutines")
	penRPS := flag.Int("penetration-rps", 1000, "Requests per second for the cache penetration test")
	penDur := flag.Duration("penetration-duration", 3*time.Second, "Duration of the cache penetration test")
	readP95 := flag.Duration("read-p95-target", 10*time.Millisecond, "Read latency p95 target")
	writeP95 := flag.Duration("write-p95-target", 50*time.Millisecond, "Write latency p95 target")
	flag.Parse()

	baseURL := strings.TrimRight(*base, "/")
	client := newHTTPClient()

	fmt.Printf("Running functional tests against %s ...\n\n", baseURL)
	funcs := runFunctionalTests(client, baseURL)

	readKey := "probe-key-000"
	for _, f := range funcs {
		if f.Key != "" {
			readKey = f.Key
			break
		}
	}

	fmt.Println("Running NFR benchmarks (this may take a while) ...")
	nfr := runBenchmarks(client, baseURL, benchmarkConfig{
		total:          *total,
		workers:        *workers,
		readKey:        readKey,
		penetrationRPS: *penRPS,
		penetrationDur: *penDur,
		readP95Target:  *readP95,
		writeP95Target: *writeP95,
	})

	if !printReport(baseURL, funcs, nfr) {
		os.Exit(1)
	}
}
