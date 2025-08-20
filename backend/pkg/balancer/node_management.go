package balancer

import (
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"os/exec"
	"strings"
	"sync"
)

// Helper method to start an internal server,
//
// In a real environment, this would not be necessary,
// and the user would just call the Balancer.AddNode method
//
// Move logic from shell script into here
func StartServer(port int, dockerInfo *config.DockerConfig) (*node.Node, error) {
	path := "./server/run.sh" //assuming you run from root of project

	cmd := exec.Command("bash", path, dockerInfo.Image, fmt.Sprintf("%d", port), fmt.Sprintf("%d", dockerInfo.InternalPort))

	output, err := cmd.Output()
	if err != nil {
		go logger.Err("Creating container", err)
		go ws.EventEmitter.Error("Creating container", err)
		return nil, err
	}
	containerID := strings.TrimSpace(string(output))
	if containerID == "" {
		err := fmt.Errorf("empty container ID received")
		go logger.Err("Creating container", err)
		go ws.EventEmitter.Error("Creating container", err)
		return nil, err
	}
	go logger.ContainerStart(containerID)
	go ws.EventEmitter.ContainerStart(containerID)

	node := node.Node{
		ContainerID:  containerID,
		Address: fmt.Sprintf("http://localhost:%d", port),
		Metrics: node.NodeMetrics{
			Health: "healthy",
		},
	}

	go logger.Info(fmt.Sprintf("Started server @ http://localhost:%d", port))
	go ws.EventEmitter.Info(fmt.Sprintf("Started server @ http://localhost:%d", port))
	return &node, nil
}

func (b *Balancer) AddNode(node *node.Node) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.Nodes = append(b.Nodes, node)
}

func (b *Balancer) RemoveNode(inputNode *node.Node) error {
	inputNode.StopServer()

	b.lock.Lock()
	defer b.lock.Unlock()

	var filtered []*node.Node
	for _, n := range b.Nodes {
		if inputNode.Equals(n) {
			filtered = append(filtered, n)
		}
	}
	b.Nodes = filtered

	return nil
}

func (b *Balancer) CleanupNodes() error {
	go logger.Info("cleaning up nodes")
	go ws.EventEmitter.Info("cleaning up nodes")
	var wg sync.WaitGroup

	for _, n := range b.Nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n.StopServer()
		}()
	}

	wg.Wait()
	b.Nodes = nil
	return nil
}
