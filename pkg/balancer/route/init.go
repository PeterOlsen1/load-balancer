package route

import (
	"fmt"
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/balancer/pool"
	"load-balancer/pkg/config"
	"load-balancer/pkg/port"
	"sync"
	"time"
)

func InitRoute(cfg config.RouteConfig) (*Route, error) {
	routeStruct := Route{
		RouteConfig: cfg,
		Queue:       InitRouteQueue(cfg.RouteQueueSize),
		NodePool:    pool.InitPool(),
	}

	var wg sync.WaitGroup

	//rethink this conditional
	if routeStruct.Docker != nil { // && len(routeStruct.Servers) == 0
		fmt.Printf("Starting containers for route: %s\n", cfg.Name)

		//start active container, don't pause so health can move to active later
		for range cfg.Pool.ActiveSize {
			wg.Add(1)
			go func() {
				defer wg.Done()

				nodePort := port.ConsumePort()
				node, err := docker.StartContainer(nodePort, routeStruct.RouteConfig)
				if err != nil {
					return
				}

				routeStruct.NodePool.AddInactive(node)
			}()
		}

		//start inactive containers and pause them
		for range cfg.Pool.InactiveSize {
			wg.Add(1)
			go func() {
				defer wg.Done()

				nodePort := port.ConsumePort()
				node, err := docker.StartContainer(nodePort, routeStruct.RouteConfig)
				if err != nil {
					return
				}

				node.Metrics.Health = "paused"
				routeStruct.NodePool.AddInactive(node)
			}()
		}
	}

	wg.Wait()
	time.Sleep(1 * time.Second)
	routeStruct.NodePool.CheckHealth(cfg)

	for _, server := range routeStruct.Servers {
		routeStruct.NodePool.AddActive(node.FromURL(server.URL, &routeStruct.RouteConfig))
	}

	//goroutine to periodically check health of containers
	go func() {
		if routeStruct.HealthTimeout <= 0 {
			//skip health check if timeout is not set
			return
		}

		//allow the server to start up before sending health request
		time.Sleep(1500 * time.Millisecond)

		ticker := time.NewTicker(time.Duration(routeStruct.HealthTimeout) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			routeStruct.NodePool.CheckHealth(cfg)
		}
	}()

	//goroutine to periodically check if we need to stop a container
	go func() {
		//allow the server to start up before sending stop requests
		time.Sleep(1500 * time.Millisecond)

		ticker := time.NewTicker(time.Duration(cfg.Pool.CleanupInterval) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			load := routeStruct.CalculateLoad()
			if load < 10 {
				routeStruct.Descale(cfg)
			}
		}
	}()

	go routeStruct.WatchQueue()

	return &routeStruct, nil
}
