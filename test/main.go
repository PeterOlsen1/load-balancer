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
const DEFAULT_SECONDS int = 0

func main() {
	numRequests := flag.Int("requests", DEFAULT_REQUESTS, "Number of requests to send")
	numSeconds := flag.Int("seconds", DEFAULT_SECONDS, "Number of seconds to send requests over")
	rps := flag.Int("rps", DEFAULT_REQUESTS, "Number of requests to send per second")
	flag.Parse()

	if *numSeconds != 0 {
		testRequestsPerSecond(*rps, *numSeconds)
	} else {
		testRequests(*numRequests, 0, true)
	}
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
			success := sendRequest("http://localhost:8080/")
			if success {
				numSuccessful++
			} else {
				numFailed++
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

func testRequestsPerSecond(rps int, seconds int) {
	fmt.Printf("Testing %d requests per second for %d seconds:\n", rps, seconds)

	var wg sync.WaitGroup
	numRequests := rps * seconds
	numSuccessful := 0
	numFailed := 0

	interval := time.Second / time.Duration(rps)
	start := time.Now()

	for i := range numRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			success := sendRequest("http://localhost:8080/")
			if success {
				numSuccessful++
			} else {
				numFailed++
			}
		}()

		nextRequestTime := start.Add(time.Duration(i+1) * interval)
		time.Sleep(time.Until(nextRequestTime))
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("\033[1m==== TESTING COMPLETE ====\033[0m")
	fmt.Printf("\033[1mTotal requests:\033[0m %d\n", numRequests)
	fmt.Printf("\033[1mSuccessful:\033[0m %d \033[1mFailed:\033[0m %d\n", numSuccessful, numFailed)
	fmt.Printf("\033[1mElapsed time:\033[0m %s\n", elapsed)
	fmt.Printf("\033[1mRequests per second:\033[0m %f\n", float64(numRequests)/elapsed.Seconds())
}

func testRequestsNTimes(numRequests int, numTests int) float64 {
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
	return avg
}

// Return true if the request was successful
func sendRequest(url string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {

		if resp != nil && resp.Body != nil {
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				fmt.Printf("Error: %v\n", readErr)
			} else {
				fmt.Printf("Error\n%s\n", string(body))
			}
		} else {
			fmt.Printf("Error: %v\n", err)
		}

		return false
	} else {
		return true
	}
}
