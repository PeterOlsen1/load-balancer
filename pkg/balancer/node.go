package balancer

import (
	"fmt"
	"load-balancer/pkg/logger"
	"net/http"
	"os/exec"
)

// add response time metric
func (node *Node) CheckHealth() error {
	address := node.Address
	resp, err := http.Get(fmt.Sprintf("%s/health", address))
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

	cmd := exec.Command("docker", "stop", node.DockerInfo.id)
	err := cmd.Run()
	if err != nil {
		go logger.LogErr("docker stop", err)
		return err
	}

	go logger.LogContainerStop(node.DockerInfo.id)
	return nil
}
