package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type benchmarkConfig struct {
	total          int
	workers        int
	readKey        string
	penetrationRPS int
	penetrationDur time.Duration
	readP95Target  time.Duration
	writeP95Target time.Duration
}

type nfrReport struct {
	TargetHost    string
	TotalRequests int
	Concurrency   int

	ReadP50       time.Duration
	ReadP95       time.Duration
	ReadP99       time.Duration
	ReadP95Target time.Duration
	ReadPass      bool

	WriteP95       time.Duration
	WriteP95Target time.Duration
	WritePass      bool

	PenetrationP95   time.Duration
	PenetrationRPS   int
	PenetrationCount int
	PenetrationPass  bool

	Successful int
	TotalRun   int
	Failed500  int
	ConnErrors int
	No500Pass  bool
}

type sample struct {
	dur    time.Duration
	status int
	err    error
}

type loadStats struct {
	durations  []time.Duration
	statuses   map[int]int
	connErrors int
	server500  int
}

func (s loadStats) total() int { return len(s.durations) }

func (s loadStats) successful() int {
	return len(s.durations) - s.connErrors - s.server500
}

func runBenchmarks(client *http.Client, baseURL string, cfg benchmarkConfig) nfrReport {
	report := nfrReport{
		TargetHost:    baseURL,
		TotalRequests: cfg.total,
		Concurrency:   cfg.workers,
	}

	readStats := runLoad(client, cfg.workers, cfg.total, readBuilder(baseURL, cfg.readKey))
	report.ReadP50 = percentile(readStats.durations, 50)
	report.ReadP95 = percentile(readStats.durations, 95)
	report.ReadP99 = percentile(readStats.durations, 99)
	report.ReadP95Target = cfg.readP95Target
	report.ReadPass = len(readStats.durations) > 0 && report.ReadP95 <= cfg.readP95Target

	report.Successful = readStats.successful()
	report.TotalRun = readStats.total()
	report.Failed500 = readStats.server500
	report.ConnErrors = readStats.connErrors

	writeStats := runLoad(client, cfg.workers, cfg.total, writeBuilder(baseURL))
	report.WriteP95 = percentile(writeStats.durations, 95)
	report.WriteP95Target = cfg.writeP95Target
	report.WritePass = len(writeStats.durations) > 0 && report.WriteP95 <= cfg.writeP95Target
	report.Failed500 += writeStats.server500
	report.ConnErrors += writeStats.connErrors

	penStats := runPaced(client, cfg.penetrationRPS, cfg.penetrationDur, cfg.workers, penetrationBuilder(baseURL))
	report.PenetrationP95 = percentile(penStats.durations, 95)
	report.PenetrationRPS = cfg.penetrationRPS
	report.PenetrationCount = penStats.total()
	report.PenetrationPass = len(penStats.durations) > 0 && report.PenetrationP95 <= cfg.readP95Target
	report.Failed500 += penStats.server500
	report.ConnErrors += penStats.connErrors

	report.No500Pass = report.Failed500 == 0

	return report
}

func runLoad(client *http.Client, workers, total int, build func(i int) (*http.Request, error)) loadStats {
	jobs := make(chan *http.Request)
	go func() {
		for i := 0; i < total; i++ {
			req, err := build(i)
			if err != nil {
				continue
			}
			jobs <- req
		}
		close(jobs)
	}()
	return runPool(client, workers, jobs)
}

func runPaced(client *http.Client, rps int, dur time.Duration, workers int, build func(i int) (*http.Request, error)) loadStats {
	interval := time.Second / time.Duration(rps)
	if interval <= 0 {
		interval = time.Microsecond
	}
	jobs := make(chan *http.Request)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		deadline := time.After(dur)
		i := 0
		for {
			select {
			case <-deadline:
				close(jobs)
				return
			case <-ticker.C:
				req, err := build(i)
				i++
				if err != nil {
					continue
				}
				jobs <- req
			}
		}
	}()
	return runPool(client, workers, jobs)
}

func runPool(client *http.Client, workers int, jobs <-chan *http.Request) loadStats {
	results := make(chan sample, workers)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range jobs {
				start := time.Now()
				resp, err := client.Do(req)
				dur := time.Since(start)
				if err != nil {
					results <- sample{dur: dur, err: err}
					continue
				}
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				results <- sample{dur: dur, status: resp.StatusCode}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	stats := loadStats{statuses: make(map[int]int)}
	for s := range results {
		stats.durations = append(stats.durations, s.dur)
		if s.err != nil {
			stats.connErrors++
			continue
		}
		stats.statuses[s.status]++
		if s.status >= http.StatusInternalServerError {
			stats.server500++
		}
	}
	return stats
}

func readBuilder(baseURL, key string) func(i int) (*http.Request, error) {
	url := baseURL + "/" + key
	return func(i int) (*http.Request, error) {
		return http.NewRequest(http.MethodGet, url, nil)
	}
}

func writeBuilder(baseURL string) func(i int) (*http.Request, error) {
	url := baseURL + "/api/v1/shorten"
	return func(i int) (*http.Request, error) {
		body := fmt.Sprintf(`{"url":"https://example.com/bench/%d-%d"}`, time.Now().UnixNano(), i)
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}
}

func penetrationBuilder(baseURL string) func(i int) (*http.Request, error) {
	return func(i int) (*http.Request, error) {
		key := fmt.Sprintf("missing-%d-%d", i, time.Now().UnixNano())
		return http.NewRequest(http.MethodGet, baseURL+"/"+key, nil)
	}
}

func percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	sorted := append([]time.Duration(nil), durations...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	idx := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
