package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const DEFAULT_REQUESTS int = 100

func main() {
	numRequests := flag.Int("requests", DEFAULT_REQUESTS, "Number of requests to send")
	flag.Parse()

	testRequests(*numRequests, 0)
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
				numFailed++

				if resp != nil && resp.Body != nil {
					body, readErr := io.ReadAll(resp.Body)
					if readErr != nil {
						fmt.Printf("Encountered error on test #%d: %v\n", i, readErr)
					} else {
						fmt.Printf("Encountered error on test #%d\n%s\n", i, string(body))
					}
				} else {
					fmt.Printf("Encountered error on test #%d: %v\n", i, err)
				}
			} else {
				fmt.Printf("Completed request #%d\n", i)
				numSuccessful++
			}
		}(i)

		if waitTime > 0 {
			time.Sleep(waitTime * time.Nanosecond)
		}
	}

	wg.Wait()
	elapsed := time.Since(start)

	avgNs := elapsed.Nanoseconds() / int64(numRequests)
	avgMs := float64(avgNs) / 1_000_000.0

	fmt.Println("\033[1m==== TESTING COMPLETE ====\033[0m")
	fmt.Printf("\033[1mAverage time per request:\033[0m %d ns (%f ms)\n", avgNs, avgMs)
	fmt.Printf("\033[1mRequests / second:\033[0m %f\n", 1000/avgMs)
	fmt.Printf("\033[1mSuccessful:\033[0m %d \033[1mFailed:\033[0m %d\n", numSuccessful, numFailed)
}
