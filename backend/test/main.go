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

	testRequests(*numRequests, 0, true)

	// testN(*numRequests, 3)
}

func testN(numRequests int, numTests int) float64 {
	res := 0.0
	for i := range numTests {
		res += 1000.0 / testRequests(numRequests, 0, false)
		time.Sleep(1 * time.Second)
		fmt.Println("Finished test", i+1)
	}
	avg := res / float64(numTests)

	fmt.Println("\033[1m==== TESTING COMPLETE ====\033[0m")
	fmt.Printf("\033[1m# of tests:\033[0m %d\n", numTests)
	fmt.Printf("\033[1mAverage req/s:\033[0m %f\n", avg)
	// fmt.Printf("\033[1mRequests / second:\033[0m %f\n", 1000/avgMs)
	// fmt.Printf("\033[1mSuccessful:\033[0m %d \033[1mFailed:\033[0m %d\n", numSuccessful, numFailed)
	return avg
}

func testRequests(numRequests int, waitTime time.Duration, log bool) float64 {
	if log {
		fmt.Printf("testing %d requests:\n", numRequests)
	}
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
				// fmt.Printf("Completed request #%d\n", i)
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

	if log {
		fmt.Println("\033[1m==== TESTING COMPLETE ====\033[0m")
		fmt.Printf("\033[1mAverage time per request:\033[0m %f ms (%d ns)\n", avgMs, avgNs)
		fmt.Printf("\033[1mRequests / second:\033[0m %f\n", 1000/avgMs)
		fmt.Printf("\033[1mSuccessful:\033[0m %d \033[1mFailed:\033[0m %d\n", numSuccessful, numFailed)
	}

	return avgMs
}
