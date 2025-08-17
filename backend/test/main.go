package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	testRequests(1000, 25*time.Millisecond)
}

func testRequests(numRequests int, duration time.Duration) {
	fmt.Printf("testing %d requests:\n", numRequests)

	start := time.Now()
	for i := range numRequests {
		_, err := http.Get("http://localhost:8080/")
		if err != nil {
			fmt.Printf("Encountered error on test #%d: %v\n", i, err)
		}
		// time.Sleep(duration)
	}
	elapsed := time.Since(start)

	avg := elapsed.Nanoseconds() / int64(numRequests)
	avgMs := avg / 1_000_000
	fmt.Printf("Average time per request: %d ns (%d ms)\n", avg, avgMs)
}
