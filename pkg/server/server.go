package server

import (
	"load-balancer/pkg/logger"
	"load-balancer/pkg/queue"
	"load-balancer/pkg/types"
	"net/http"
)

func Serve() error {
	http.HandleFunc("/", connectionHandler)
	return http.ListenAndServe(":8080", nil)
}

func connectionHandler(resp http.ResponseWriter, req *http.Request) {
	conn := types.Connection{
		Writer:  resp,
		Request: req,
	}

	go logger.LogRequest(&conn)
	queue.ConnectionQueue.Enqueue(&conn)
}
