package main

// var wg sync.WaitGroup

import (
	"flag"
	"fmt"
	"load-balancer/pkg/balancer"
	_ "load-balancer/pkg/balancer/receiver"
	"load-balancer/pkg/config"
	"load-balancer/pkg/server"
	"os"
	"os/signal"
)

func main() {
	configPath := flag.String("cfg", "./config/config.yaml", "Location of configuration file")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		os.Exit(2)
	}

	fmt.Println(config.Config)

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
	balancer.LoadBalancer.InitBalancer(0)
	server.Serve(config.Config.Server.Host, config.Config.Server.Port)
}
