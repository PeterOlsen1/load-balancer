package balancer

import (
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"sync"
)

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
