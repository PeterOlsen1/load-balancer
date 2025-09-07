package route

import (
	"load-balancer/pkg/balancer/docker"
	"sync"
)

func (r *Route) CleanupNodes() error {
	var wg sync.WaitGroup
	var loopErr error

	// capture lock so no other processes can add containers
	r.NodePool.Mu.Lock()
	defer r.NodePool.Mu.Unlock()

	for _, n := range r.NodePool.GetAll() {
		if n == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			n.CloseQueue()
			err := docker.StopContainer(n.ContainerID)
			if err != nil {
				loopErr = err
			}
		}()
	}

	wg.Wait()
	return loopErr
}
