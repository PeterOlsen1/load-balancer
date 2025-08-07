package server

import (
	"fmt"
	"load-balancer/pkg/types"
	"load-balancer/pkg/queue"
	"net/http"
)

func Serve() error {
	http.HandleFunc("/", connectionHandler)
	return http.ListenAndServe(":8080", nil)
}

func connectionHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Got request to:")
	fmt.Println(req.URL)

	conn := types.Connection{
		Writer:  resp,
		Request: req,
	}

	queue.ConnectionQueue.Push(&conn)
}
