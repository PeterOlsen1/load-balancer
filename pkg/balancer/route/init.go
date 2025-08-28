package route

import (
	"load-balancer/pkg/balancer/docker"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/balancer/pool"
	"load-balancer/pkg/config"
	"load-balancer/pkg/port"
	"time"
)

func InitRoute(cfg config.RouteConfig) (*Route, error) {
	routeStruct := Route{
		RouteConfig: cfg,
		Queue:       InitRouteQueue(),
		NodePool:    pool.InitPool(),
	}

	//rethink this conditional
	if routeStruct.Docker != nil && len(routeStruct.Servers) == 0 {
		//start # initial docker containers, add to inactive pool
		for range cfg.Docker.InitialContainers {
			port := port.ConsumePort()
			node, err := docker.StartContainer(port, routeStruct.RouteConfig)
			if err != nil {
				return nil, err
			}

			routeStruct.NodePool.AddInactive(node)
		}

		//call the scale method here to refill the inactive pool
		routeStruct.Scale(routeStruct.RouteConfig)
	}

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
			routeStruct.NodePool.CheckHealth()
		}
	}()

	//goroutine to periodically check if we need to stop a container
	go func() {
		if routeStruct.InactiveTimeout <= 0 {
			//this might be a bad idea but I'm not sure how a negative time would work anyway
			return
		}

		ticker := time.NewTicker(time.Duration(routeStruct.HealthTimeout) * time.Millisecond)
		defer ticker.Stop()

		// for range ticker.C {

		// }
	}()

	go routeStruct.WatchQueue()

	return &routeStruct, nil
}
