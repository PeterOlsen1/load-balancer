package server

import (
	"fmt"
	"load-balancer/pkg/queue"
	"load-balancer/pkg/types"
	"net/http"
)

func Serve() error {
	http.HandleFunc("/", connectionHandler)
	return http.ListenAndServe(":8080", nil)
}

func connectionHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Adding request to queue")

	conn := types.Connection{
		Writer:  resp,
		Request: req,
	}

	queue.ConnectionQueue.Enqueue(&conn)
}
