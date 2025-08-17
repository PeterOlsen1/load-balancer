package server

import (
	"fmt"
	"load-balancer/pkg/balancer"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"load-balancer/pkg/ws"
	"net/http"
)

func Serve(address string, port int) error {
	url := fmt.Sprintf("%s:%d", address, port)
	fmt.Println("Starting server @", url)
	http.HandleFunc("/ws", ws.WsHandler)
	http.HandleFunc("/new", addNewContainer)
	http.HandleFunc("/", connectionHandler)
	return http.ListenAndServe(url, nil)
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
	node, err := balancer.StartServer(balancer.PORT)
	if err != nil {
		return
	}

	balancer.LoadBalancer.AddNode(node)
	balancer.PORT++
	fmt.Fprintf(resp, "Added new container: %s", node.DockerInfo.Id)
}
