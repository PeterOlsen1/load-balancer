package main

import (
	"flag"
	"fmt"
	"load-balancer/pkg/balancer"
	_ "load-balancer/pkg/balancer/receiver"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/server"
	"os"
	"os/signal"
)

func main() {
	configPath := flag.String("cfg", "./config/config.yaml", "Location of configuration file")
	flag.Parse()

	err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Could not read config! Exiting...")
		os.Exit(1)
	}
	logger.InitLogger()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		err := balancer.Balancer.CleanupNodes()
		if err != nil {
			fmt.Println("Error cleaning up nodes:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	//unused
	// go balancer.WatchQueue()
	balancer.Balancer.InitBalancer()
	server.Serve(config.Config.Server.Host, config.Config.Server.Port)
}
