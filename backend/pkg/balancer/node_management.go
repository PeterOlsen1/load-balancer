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

	fmt.Println(path, dockerInfo.Image, port, dockerInfo.InternalPort)
	port := ConsumePort()
	cmd := exec.Command("bash", path, dockerInfo.Image, fmt.Sprintf("%d", port), fmt.Sprintf("%d", dockerInfo.InternalPort))

	output, err := cmd.Output()
	if err != nil {
		logger.Err("Creating container", err)
		ws.EventEmitter.Error("Creating container", err)
		return nil, err
	}
	containerID := strings.TrimSpace(string(output))
	if containerID == "" {
		err := fmt.Errorf("empty container ID received")
		logger.Err("Creating container", err)
		ws.EventEmitter.Error("Creating container", err)
		return nil, err
	}
	logger.ContainerStart(containerID)
	ws.EventEmitter.ContainerStart(containerID)

	node := node.Node{
		ContainerID: containerID,
		Address:     fmt.Sprintf("http://localhost:%d", port),
		Metrics: node.NodeMetrics{
			Health: "healthy",
		},
	}

	logger.Info(fmt.Sprintf("Started server @ http://localhost:%d", port))
	ws.EventEmitter.Info(fmt.Sprintf("Started server @ http://localhost:%d", port))
	return &node, nil
}

func (r *Route) AddNode(inputNode *node.Node) {
	Balancer.NodeTable[inputNode.ContainerID] = inputNode

	r.lock.Lock()
	defer r.lock.Unlock()
	r.Nodes = append(r.Nodes, inputNode)
}

// No lock since the only place RemoveNode is called already
// has the Route lock acquired
func (r *Route) RemoveNode(inputNode *node.Node) error {

	var filtered []*node.Node
	for _, n := range r.Nodes {
		if !inputNode.Equals(n) {
			filtered = append(filtered, n)
		}
	}

	err := inputNode.StopServer()
	if err != nil {
		return err
	}

	r.Nodes = filtered
	delete(Balancer.NodeTable, inputNode.ContainerID)

	for _, n := range r.Nodes {
		fmt.Println(n.Address)
	}
	return nil
}

func (r *Route) CleanupNodes() error {
	var wg sync.WaitGroup
	var loopErr error

	for _, n := range r.Nodes {
		if n == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := n.StopServer()
			if err != nil {
				loopErr = err
			}
		}()
	}

	wg.Wait()
	r.Nodes = nil
	return loopErr
}

func (b *BalancerType) CleanupNodes() error {
	logger.Info("cleaning up nodes")
	ws.EventEmitter.Info("cleaning up nodes")
	var wg sync.WaitGroup
	var loopErr error

	for _, r := range b.Routes {
		if r == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.CleanupNodes()
			if err != nil {
				loopErr = err
			}
		}()
	}

	wg.Wait()
	return loopErr
}
