package node

import "load-balancer/pkg/config"

func FromContainer(containerID string, address string, routeConfig *config.RouteConfig) *Node {
	out := &Node{
		ContainerID: containerID,
		Address:     address,
		Metrics: NodeMetrics{
			Health:       "unknown",
			ResponseTime: 0,
			Connections:  0,
		},
		Queue: *InitNodeQueue(routeConfig.RequestLimit),
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
		Queue: *InitNodeQueue(routeConfig.RequestLimit),
	}

	go out.CheckHealth()
	go out.WatchQueue()
	return out
}
