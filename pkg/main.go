package main

import (
	"load-balancer/pkg/balancer"
	"load-balancer/pkg/server"
)

// var wg sync.WaitGroup

func main() {

	go balancer.WatchQueue()
	server.Serve()
}
