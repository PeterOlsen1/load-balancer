package node

import (
	"load-balancer/pkg/config"
	"net/http"
	"time"
)

var transport = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 30,
	IdleConnTimeout:     30 * time.Second,
}

var httpClient = &http.Client{
	Transport: transport,
	Timeout:   10 * time.Second,
}

func FromContainer(containerID string, address string, routeConfig config.RouteConfig) *Node {
	out := &Node{
		ContainerID: containerID,
		Address:     address,
		Metrics: NodeMetrics{
			Health:       "unknown",
			ResponseTime: 0,
			Connections:  0,
		},
		Queue: InitNodeQueue(routeConfig.NodeQueueSize),
	}

	go out.CheckHealth()
	go out.WatchQueue()
	return out
}

func FromURL(url string, routeConfig *config.RouteConfig) *Node {
	out := &Node{
		ContainerID: "",
		Address:     url,
		Metrics: NodeMetrics{
			Health:       "unknown",
			ResponseTime: 0,
			Connections:  0,
		},
		Queue: InitNodeQueue(routeConfig.NodeQueueSize),
	}

	go out.CheckHealth()
	go out.WatchQueue()
	return out
}
