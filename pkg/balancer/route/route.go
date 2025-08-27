package route

import (
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/port"
	"sync"
)

// the Scale() method should automatically decide whether to spin
// up a docker container, or use a pre-existing URL.
//
// The goal here is that we'll have a few containers to
// pick from, if we use one, make sure to warm up another
func (r *Route) Scale() (*node.Node, error) {
	port := port.ConsumePort()
	node, err := docker.StartContainer(port, &r.RouteConfig)
	if err != nil {
		return nil, err
	}

	r.NodePool.AddInactive(node)
	return node, nil
}

func (r *Route) CleanupNodes() error {
	var wg sync.WaitGroup
	var loopErr error

	for _, n := range r.NodePool.GetAll() {
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
	return loopErr
}
