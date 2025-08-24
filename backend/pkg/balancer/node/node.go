package node

import (
	"fmt"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"net/http"
	"time"
)

func FromURL(url string) *Node {
	out := &Node{
		ContainerID: "",
		Address:     url,
		Metrics: NodeMetrics{
			Health:       "unknown",
			ResponseTime: 0,
			Connections:  0,
		},
	}

	go out.CheckHealth()
	return out
}

// Send a request to the node backend to check the health
//
// If an OK status is returned, set node to healthy. Else, unhealthy
func (node *Node) CheckHealth() error {
	node.Metrics.Lock.Lock()
	isPaused := node.Metrics.Health == "paused"
	node.Metrics.Lock.Unlock()

	if isPaused {
		return nil
	}

	address := node.Address

	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/health", address))
	duration := time.Since(start)

	if err != nil {
		logger.Err("Fetching node health", err)
		ws.EventEmitter.Error("Fetching node health", err)
		return err
	}

	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()
	
	respTime := float32(duration.Microseconds() / 1000)
	node.Metrics.ResponseTime = respTime

	health := "healthy"
	if resp.StatusCode != http.StatusOK {
		health = "unhealthy"
		go node.Queue.CloseQueue()

		logger.Health(health, node.Address, respTime)
		ws.EventEmitter.Health(health, node.Address, respTime)
	} else {
		logger.Health(health, node.Address, respTime)
		ws.EventEmitter.Health(health, node.Address, respTime)
	}
	node.Metrics.Health = health

	return nil
}

func (node *Node) Pause() {
	node.Metrics.Lock.Lock()
	node.Metrics.Health = "paused"
	node.Metrics.Lock.Unlock()

	node.Queue.CloseQueue()
}

func (node *Node) Unpause() {
	node.Metrics.Lock.Lock()
	node.Metrics.Health = "unknown"
	node.Metrics.Lock.Unlock()

	node.CheckHealth()
}

func (n *Node) Equals(other *Node) bool {
	return n.Address == other.Address && n.ContainerID == other.ContainerID
}
