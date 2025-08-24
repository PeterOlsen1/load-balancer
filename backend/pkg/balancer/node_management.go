package balancer

import (
	"fmt"
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"sync"
)

func (r *Route) Scale() (*node.Node, error) {
	port := ConsumePort()
	return docker.StartContainer(port, r.Docker)
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
