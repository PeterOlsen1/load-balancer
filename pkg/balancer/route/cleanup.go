package route

import (
	"load-balancer/pkg/balancer/docker"
	"sync"
)

func (r *Route) CleanupNodes() error {
	var wg sync.WaitGroup
	var loopErr error

	// locks mutex and does not unlock, health operations are also stopped
	r.NodePool.Close()
	for _, n := range r.NodePool.GetAll() {
		if n == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			n.CloseQueue()

			// URL node
			if n.ContainerID == "" {
				return
			}

			err := docker.StopContainer(n.ContainerID)
			if err != nil {
				loopErr = err
			}
		}()
	}

	wg.Wait()
	return loopErr
}
