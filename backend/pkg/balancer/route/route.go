package route

import (
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/port"
	"sync"
)

// the Scale() method should automatically decide whether to spin
// up a docker container, or use a pre-existing URL.
func (r *Route) Scale() (*node.Node, error) {
	port := port.ConsumePort()
	node, err := docker.StartContainer(port, &r.RouteConfig)
	if err != nil {
		return nil, err
	}

	r.addNode(node)
	return node, nil
}

func (r *Route) addNode(inputNode *node.Node) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
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

	err := docker.StopContainer(inputNode.ContainerID)
	if err != nil {
		return err
	}

	r.Nodes = filtered
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
			// n.CloseQueue()
			err := docker.StopContainer(n.ContainerID)
			if err != nil {
				loopErr = err
			}
		}()
	}

	wg.Wait()
	r.Nodes = nil
	return loopErr
}
