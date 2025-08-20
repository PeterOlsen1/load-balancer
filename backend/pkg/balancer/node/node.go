package node

import (
	"fmt"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"net/http"
	"os/exec"
	"time"
)

func FromURL(url string) *Node {
	return &Node{
		ContainerID: "",
		Address:     url,
		Metrics: NodeMetrics{
			Health:       "unknown",
			ResponseTime: 0,
			Connections:  0,
		},
	}
}

// Send a request to the node backend to check the health
//
// If an OK status is returned, set node to healthy. Else, unhealthy
func (node *Node) CheckHealth() error {
	node.Metrics.Lock.Lock()
	if node.Metrics.Health == "paused" {
		return nil
	}
	node.Metrics.Lock.Unlock()
	address := node.Address

	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/health", address))
	duration := time.Since(start)

	if err != nil {
		go logger.Err("Fetching node health", err)
		go ws.EventEmitter.Error("Fetching node health", err)
		return err
	}
	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()
	respTime := float32(duration.Microseconds() / 1000)
	node.Metrics.ResponseTime = respTime

	health := "healthy"
	if resp.StatusCode != http.StatusOK {
		health = "unhealthy"
		go logger.Health(health, node.Address, respTime)
		go ws.EventEmitter.Health(health, node.Address, respTime)
	} else {
		go logger.Health(health, node.Address, respTime)
		go ws.EventEmitter.Health(health, node.Address, respTime)
	}
	node.Metrics.Health = health

	return nil
}

func (node *Node) Pause() {
	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()

	node.Metrics.Health = "paused"
}

func (node *Node) Unpause() {
	node.Metrics.Lock.Lock()
	node.Metrics.Health = "unknown"
	node.Metrics.Lock.Unlock()

	node.CheckHealth()
}

// Stops the server associated with any given node
// through the docker stop command.
//
// If this node has no server, instantly return nil
func (node *Node) StopServer() error {
	if node.ContainerID == "" {
		return nil
	}

	cmd := exec.Command("docker", "stop", node.ContainerID)
	err := cmd.Run()
	if err != nil {
		go logger.Err("docker stop", err)
		go ws.EventEmitter.Error("docker stop", err)
		return err
	}

	go logger.ContainerStop(node.ContainerID)
	go ws.EventEmitter.ContainerStop(node.ContainerID)
	return nil
}

func (n *Node) Equals(other *Node) bool {
	return n.Address == other.Address && n.ContainerID == other.ContainerID
}

// Returns a node from a URL, instead of spinning up a docker container.
// This is to be used when the user already has a service running,
// and wants to just input it as a node.
//
// This would require interaction from the frontend
func FromUrl(url string) *Node {
	out := Node{
		Address: url,
	}

	go out.CheckHealth()
	return &out
}

/*
	f, err := os.OpenFile("./data/urls", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Err("Failed to open url file", err)
		return nil, err
	}
	defer f.Close()
	if _, err := f.WriteString(url + "\n"); err != nil {
		logger.Err("Failed to write to url file", err)
		return nil, err
	}
*/
