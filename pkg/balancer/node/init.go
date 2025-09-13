package node

import (
	"load-balancer/pkg/config"
	"load-balancer/pkg/types"
	"load-balancer/pkg/workerpool"
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

func (n *Node) getWorkerPoolEventHandler() func(*types.Connection) {
	return func(conn *types.Connection) {
		n.processRequest(conn)
	}
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
	}
	// add node queue later since we need to call the n.getWorkerPoolEventHandler method
	out.Queue = InitNodeQueue(routeConfig.RouteQueueSize, routeConfig.WorkerThreads, out.getWorkerPoolEventHandler())

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
	}
	// add node queue later since we need to call the n.getWorkerPoolEventHandler method
	out.Queue = InitNodeQueue(routeConfig.NodeQueueSize, routeConfig.WorkerThreads, out.getWorkerPoolEventHandler())

	go out.CheckHealth()
	go out.WatchQueue()
	return out
}

func InitNodeQueue(capacity uint32, workerThreads uint16, eventHandler func(*types.Connection)) *NodeQueue {
	return &NodeQueue{
		queue:       make(chan *types.Connection, capacity),
		open:        true,
		connChan:    make(chan *types.Connection, capacity),
		closeSignal: make(chan struct{}),
		workerPool:  workerpool.InitWorkerPool(workerThreads, eventHandler),
	}
}
