package node

import (
	"fmt"
	"load-balancer/pkg/logger"
	"net/http"
	"os/exec"
	"time"
)

func (node *Node) CheckHealth() error {
	address := node.Address

	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/health", address))
	duration := time.Since(start)

	if err != nil {
		go logger.LogErr("Fetching node health", err)
		return err
	}

	health := Healthy
	if resp.StatusCode != http.StatusOK {
		health = Unhealthy
		go logger.LogStatusCheck("Unhealthy", node.Address)
	} else {
		go logger.LogStatusCheck("Healthy", node.Address)
	}
	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()
	node.Metrics.Health = health
	node.Metrics.ResponseTime = float32(duration.Microseconds() / 1000)

	return nil
}

// Stops the server associated with any given node
// through the docker stop command.
//
// If this node has no server, instantly return nil
func (node *Node) StopServer() error {
	if node.DockerInfo == nil {
		return nil
	}

	cmd := exec.Command("docker", "stop", node.DockerInfo.Id)
	err := cmd.Run()
	if err != nil {
		go logger.LogErr("docker stop", err)
		return err
	}

	go logger.LogContainerStop(node.DockerInfo.Id)
	return nil
}

func (node *Node) GetHealth() string {
	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()

	health := node.Metrics.Health
	switch health {
	case 0:
		return "Unknown"
	case 1:
		return "Unhealthy"
	default:
		return "Healthy"
	}
}

func (n *Node) Equals(other *Node) bool {
	return n.Address == other.Address && n.DockerInfo.Id == other.DockerInfo.Id
}
