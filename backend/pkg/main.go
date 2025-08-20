package main

// var wg sync.WaitGroup

import (
	"flag"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		balancer.Balancer.CleanupNodes()
	}()

	//unused
	// go balancer.WatchQueue()
	balancer.Balancer.InitBalancer(0)
	server.Serve(config.Config.Server.Host, config.Config.Server.Port)
}
