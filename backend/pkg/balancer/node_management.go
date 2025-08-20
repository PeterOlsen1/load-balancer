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
func StartServer(dockerInfo *config.DockerConfig) (*node.Node, error) {
	path := "./server/run.sh" //assuming you run from root of project

	port := ConsumePort()
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
		ContainerID: containerID,
		Address:     fmt.Sprintf("http://localhost:%d", port),
		Metrics: node.NodeMetrics{
			Health: "healthy",
		},
	}

	go logger.Info(fmt.Sprintf("Started server @ http://localhost:%d", port))
	go ws.EventEmitter.Info(fmt.Sprintf("Started server @ http://localhost:%d", port))
	return &node, nil
}

func (r *Route) AddNode(node *node.Node) {
	Balancer.NodeTable[node.ContainerID] = node

	r.lock.Lock()
	defer r.lock.Unlock()

	r.Nodes = append(r.Nodes, node)
}

func (r *Route) RemoveNode(inputNode *node.Node) error {
	inputNode.StopServer()

	r.lock.Lock()
	defer r.lock.Unlock()

	var filtered []*node.Node
	for _, n := range r.Nodes {
		if inputNode.Equals(n) {
			filtered = append(filtered, n)
		}
	}
	r.Nodes = filtered

	delete(Balancer.NodeTable, inputNode.ContainerID)

	return nil
}

func (r *Route) CleanupNodes() {
	var wg sync.WaitGroup

	for _, n := range r.Nodes {
		if n == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			n.StopServer()
		}()
	}

	wg.Wait()
	r.Nodes = nil
}

func (b *BalancerType) CleanupNodes() {
	go logger.Info("cleaning up nodes")
	go ws.EventEmitter.Info("cleaning up nodes")
	var wg sync.WaitGroup

	for _, r := range b.Routes {
		if r == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			r.CleanupNodes()
		}()
	}

	wg.Wait()
}
