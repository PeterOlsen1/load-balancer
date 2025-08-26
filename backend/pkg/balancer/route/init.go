package route

import (
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"time"
)

func InitRoute(cfg config.RouteConfig) (*Route, error) {
	routeStruct := Route{
		RouteConfig: cfg,
		Queue:       InitRouteQueue(),
	}

	if routeStruct.Docker != nil && len(routeStruct.Servers) == 0 {
		for i := range cfg.Docker.InitialContainers {
			_, err := routeStruct.Scale()
			if err != nil {
				fmt.Printf("failed starting initial container #%d\n", i)
				return nil, err
			}
		}
	}

	for _, server := range routeStruct.Servers {
		routeStruct.Nodes = append(routeStruct.Nodes, node.FromURL(server.URL, &routeStruct.RouteConfig))
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
			routeStruct.Lock.Lock()
			for _, n := range routeStruct.Nodes {
				go n.CheckHealth()
			}
			routeStruct.Lock.Unlock()
		}
	}()

	// //goroutine to periodically check if we need to stop a container
	// go func() {
	// 	if routeStruct.InactiveTimeout <= 0 {
	// 		//this might be a bad idea but I'm not sure how a negative time would work anyway
	// 		return
	// 	}

	// 	ticker := time.NewTicker(time.Duration(routeStruct.HealthTimeout) * time.Millisecond)
	// 	defer ticker.Stop()

	// 	for range ticker.C {
	// 		routeStruct.Lock.Lock()
	// 		for i := len(routeStruct.Nodes) - 1; i >= 0; i-- {
	// 			node := routeStruct.Nodes[i]
	// 			node.Metrics.Lock.Lock()
	// 			if routeStruct.Docker != nil && len(routeStruct.Nodes) > 1 && time.Since(node.Metrics.LastRequestTime).Milliseconds() > time.Duration(routeStruct.InactiveTimeout).Milliseconds() {
	// 				delete(Balancer.NodeTable, node.ContainerID)
	// 				routeStruct.RemoveNode(node)
	// 			}
	// 			node.Metrics.Lock.Unlock()
	// 		}
	// 		routeStruct.Lock.Unlock()
	// 	}
	// }()

	go routeStruct.WatchQueue()

	return &routeStruct, nil
}
