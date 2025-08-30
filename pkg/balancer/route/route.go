package route

import (
	"fmt"
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/port"
	"sync"
	"time"
)

// the Scale() method should automatically decide whether to spin
// up a docker container, or use a pre-existing URL.
//
// The goal here is that we'll have a few containers to
// pick from, if we use one, make sure to warm up another
func (r *Route) Scale(cfg config.RouteConfig) error {
	if time.Since(r.LastScale) < 500*time.Millisecond {
		return nil
	}

	r.LastScale = time.Now()
	err := r.NodePool.UnpauseOne()

	//err will != nil when len(inactive) == 0
	if err != nil {
		logger.Info(fmt.Sprintf("zero inactive containers, adding %d", cfg.Pool.InactiveSize))
		for range cfg.Pool.InactiveSize {
			port := port.ConsumePort()
			node, err := docker.StartContainer(port, r.RouteConfig)
			if err != nil {
				return err
			}

			node.Metrics.Health = "paused"
			r.NodePool.AddInactive(node)
		}
	} else if r.NodePool.GetInactiveSize() < cfg.Pool.InactiveSize {
		//always keep cfg.Docker.InitialContainers in the inactive pool
		logger.Info(fmt.Sprintf("fewer inactive nodes than initial docker containers, adding %d", cfg.Pool.InactiveSize-r.NodePool.GetInactiveSize()))

		for range cfg.Pool.InactiveSize - r.NodePool.GetInactiveSize() {
			port := port.ConsumePort()
			node, err := docker.StartContainer(port, r.RouteConfig)
			if err != nil {
				return err
			}

			r.NodePool.AddInactive(node)
			node.Metrics.Health = "paused"
		}
	}

	fmt.Println("Node pools after scale")
	fmt.Println(r.NodePool.Active)
	fmt.Println(r.NodePool.Inactive)

	return nil
}

// Scale down the amount of containers we have running only
// if there are more than the initial amount
func (r *Route) Descale(cfg config.RouteConfig) {
	if r.NodePool.GetActiveSize() > cfg.Pool.ActiveSize {
		fmt.Println("Descaling...")
		err := r.NodePool.PauseOne()
		if err != nil {
			logger.Err("descaling one container", err)
		}
	}
}

func (r *Route) CalculateLoad() float64 {
	conns := r.Queue.Len()
	numNodes := r.NodePool.GetActiveSize()
	maxCapacity := numNodes * r.RouteConfig.RequestLimit

	if maxCapacity <= 0 {
		return 0
	}

	for _, n := range r.NodePool.GetActive() {
		conns += n.Queue.Len()
	}

	return (float64(conns) / float64(maxCapacity)) * 100
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
