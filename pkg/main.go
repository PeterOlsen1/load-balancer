package main

import (
	"load-balancer/pkg/balancer"
	"load-balancer/pkg/server"
	"os"
	"os/signal"
)

// var wg sync.WaitGroup

func main() {

	balancer.LoadBalancer.InitBalancer()
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

	go balancer.WatchQueue()
	server.Serve()
}
