package main

// var wg sync.WaitGroup

import (
	"flag"
	"load-balancer/pkg/balancer"
	"load-balancer/pkg/server"
	"os"
	"os/signal"
)

func main() {
	address := flag.String("addr", "127.0.0.1", "Address to run the server on")
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		err := balancer.LoadBalancer.CleanupNodes()
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	//unused
	// go balancer.WatchQueue()
	// balancer.LoadBalancer.InitBalancer()
	server.Serve(*address, *port)
}
