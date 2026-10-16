package main

import (
	"fmt"
	"time"
)

func printReport(baseURL string, funcs []functionalResult, nfr nfrReport) bool {
	allFuncPass := true
	for _, f := range funcs {
		if !f.Passed {
			allFuncPass = false
		}
	}

	fmt.Println()
	fmt.Println("==================================================")
	fmt.Println("              URL SHORTENER JUDGE REPORT          ")
	fmt.Println("==================================================")
	fmt.Println("[FUNCTIONAL TESTS]")
	for _, f := range funcs {
		status := "PASS"
		if !f.Passed {
			status = "FAIL"
		}
		fmt.Printf("  [%s] %s\n", status, f.Name)
		if f.Detail != "" {
			fmt.Printf("         -> %s\n", f.Detail)
		}
	}

	fmt.Println()
	fmt.Println("[NFR & BENCHMARK RESULTS]")
	fmt.Printf("  Target Host         : %s\n", nfr.TargetHost)
	fmt.Printf("  Total Requests      : %s\n", commas(nfr.TotalRequests))
	fmt.Printf("  Concurrency Level   : %d Workers\n", nfr.Concurrency)
	fmt.Println()
	fmt.Printf("  Read Latency (p50)  : %s\n", ms(nfr.ReadP50))
	fmt.Printf("  Read Latency (p95)  : %s  [%s %s %s]\n",
		ms(nfr.ReadP95), verdictTag(nfr.ReadPass), cmp(nfr.ReadPass), ms(nfr.ReadP95Target))
	fmt.Printf("  Read Latency (p99)  : %s\n", ms(nfr.ReadP99))
	fmt.Println()
	fmt.Printf("  Write Latency (p95) : %s [%s %s %s]\n",
		ms(nfr.WriteP95), verdictTag(nfr.WritePass), cmp(nfr.WritePass), ms(nfr.WriteP95Target))
	fmt.Printf("  Cache Penetration   : %s  [%s]  (%s reqs @ %d rps)\n",
		ms(nfr.PenetrationP95), verdictTag(nfr.PenetrationPass), commas(nfr.PenetrationCount), nfr.PenetrationRPS)
	fmt.Printf("  Successful Requests : %s / %s (%.0f%%)\n",
		commas(nfr.Successful), commas(nfr.TotalRun), percent(nfr.Successful, nfr.TotalRun))
	fmt.Printf("  Failed / 500 Errors : %d\n", nfr.Failed500)
	if nfr.ConnErrors > 0 {
		fmt.Printf("  Connection Errors   : %d\n", nfr.ConnErrors)
	}

	passed := allFuncPass && nfr.ReadPass && nfr.WritePass && nfr.PenetrationPass && nfr.No500Pass

	fmt.Println()
	fmt.Printf("FINAL VERDICT: %s\n", verdictWord(passed))
	fmt.Println("==================================================")

	return passed
}

func verdictWord(ok bool) string {
	if ok {
		return "PASSED"
	}
	return "FAILED"
}

func verdictTag(ok bool) string {
	if ok {
		return "PASS"
	}
	return "FAIL"
}

func cmp(ok bool) string {
	if ok {
		return "<="
	}
	return ">"
}

func ms(d time.Duration) string {
	return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000.0)
}

func percent(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func commas(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	return string(out)
}
