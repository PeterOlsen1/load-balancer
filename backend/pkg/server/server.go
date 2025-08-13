package server

import (
	"fmt"
	"load-balancer/pkg/balancer"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"load-balancer/pkg/ws"
	"net/http"
)

var PORT int = 3001

func Serve() error {
	fmt.Println("Starting server @ localost:8080")
	http.HandleFunc("/ws", ws.WsHandler)
	http.HandleFunc("/", connectionHandler)
	http.HandleFunc("/new", addNewContainer)
	return http.ListenAndServe(":8080", nil)
}

func connectionHandler(resp http.ResponseWriter, req *http.Request) {
	conn := types.Connection{
		Response: resp,
		Request:  req,
	}

	fmt.Println(req.Method + ": " + req.URL.Path)
	go logger.Request(&conn)
	go ws.EventEmitter.Request(&conn)
	balancer.LoadBalancer.ProxyRequest(&conn)
}

// test endpoint for adding new container functionality
func addNewContainer(resp http.ResponseWriter, req *http.Request) {
	node, err := balancer.StartServer(PORT)
	if err != nil {
		return
	}

	balancer.LoadBalancer.AddNode(node)
	PORT++
	fmt.Fprintf(resp, "Added new container: %s", node.DockerInfo.Id)
}
