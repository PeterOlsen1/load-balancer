package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	testRequests(50, 0)
}

func testRequests(numRequests int, waitTime time.Duration) {
	fmt.Printf("testing %d requests:\n", numRequests)
	var wg sync.WaitGroup

	numFailed := 0
	numSuccessful := 0

	start := time.Now()
	for i := range numRequests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080/")
			if err != nil || resp.StatusCode != 200 {
				fmt.Printf("Encountered error on test #%d: %v\n", i, err)
				numFailed++
			} else {
				fmt.Printf("Completed request #%d\n", i)
				numSuccessful++
			}
		}(i)

		time.Sleep(waitTime * time.Millisecond)
	}

	wg.Wait()
	elapsed := time.Since(start)

	avg := elapsed.Nanoseconds() / int64(numRequests)
	avgMs := avg / 1_000_000
	fmt.Printf("Average time per request: %d ns (%d ms)\n", avg, avgMs)
	fmt.Printf("Successful: %d Failed: %d\n", numSuccessful, numFailed)
}
