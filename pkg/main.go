package main

import (
	"load-balancer/pkg/server"
	// "sync"
)

// var wg sync.WaitGroup

func main() {

	server.Serve()
}

// func httpTest() (*string, error) {
// 	defer wg.Done()

// 	resp, err := http.Get("http://example.com")
// 	if err != nil {
// 		fmt.Println("there was an error fetching example!", err)
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("there was an reading fetching example!", err)
// 		return nil, err
// 	}

// 	strBody := string(body)
// 	fmt.Println("response:", strBody)
// 	return &strBody, nil
// }
