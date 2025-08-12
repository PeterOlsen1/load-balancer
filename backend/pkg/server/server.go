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
	balancer.LoadBalancer.ProxyRequest(&conn)
}
