package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	testRequests(10000)
}

func testRequests(numRequests int) {
	fmt.Printf("testing %d requests:\n", numRequests)
	var wg sync.WaitGroup

	start := time.Now()
	for i := range numRequests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := http.Get("http://localhost:8080/")
			if err != nil {
				fmt.Printf("Encountered error on test #%d: %v\n", i, err)
			}
		}(i)
	}
	elapsed := time.Since(start)

	avg := elapsed.Nanoseconds() / int64(numRequests)
	avgMs := avg / 1_000_000
	fmt.Printf("Average time per request: %d ns (%d ms)\n", avg, avgMs)
}
