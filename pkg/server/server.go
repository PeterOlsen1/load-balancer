package server

import (
	"fmt"
	"load-balancer/pkg/balancer"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"net/http"
)

func Serve() error {
	fmt.Println("Starting server @ localost:8080")
	http.HandleFunc("/", connectionHandler)
	return http.ListenAndServe(":8080", nil)
}

func connectionHandler(resp http.ResponseWriter, req *http.Request) {
	conn := types.Connection{
		Response: resp,
		Request:  req,
	}

	fmt.Println(req.Method + ": " + req.URL.Path)
	go logger.LogRequest(&conn)
	// queue.ConnectionQueue.Enqueue(&conn)
	balancer.LoadBalancer.ProxyRequest(&conn)
}

func Send500(conn *types.Connection) {
	message := "500 Internal Server Error"
	conn.Response.Header().Set("Content-Type", "text/plain")
	conn.Response.Header().Set("Content-Length", fmt.Sprintf("%d", len(message)))
	conn.Response.WriteHeader(500)

	// Write the response body once
	_, err := conn.Response.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing 500 response:", err)
	}
}
