package balancer

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/balancer/route"
	"load-balancer/pkg/config"
	"load-balancer/pkg/queue"
	"time"
)

var Balancer = BalancerType{
	NodeTable: make(map[string]*node.Node),
}

func WatchQueue() {
	for {
		<-queue.ConnectionQueue.Notify
		for {
			conn, err := queue.ConnectionQueue.Dequeue()

			if err != nil {
				break
			}

			Balancer.ProxyRequest(conn)
		}
	}
}

// Pass in num <= 0 for no health checks
func (b *BalancerType) InitBalancer() error {
	for _, r := range config.Config.Routes {
		routeStruct := route.Route{
			RouteConfig: r,
		}

		b.Routes = append(b.Routes, &routeStruct)
		if routeStruct.Docker != nil && len(routeStruct.Servers) == 0 {
			serverNode, err := routeStruct.Scale()
			if err != nil {
				return err
			}

			routeStruct.Lock.Lock()
			routeStruct.Nodes = append(routeStruct.Nodes, serverNode)
			routeStruct.Lock.Unlock()
			b.NodeTable[serverNode.ContainerID] = serverNode
		}

		for _, server := range routeStruct.Servers {
			routeStruct.Nodes = append(routeStruct.Nodes, node.FromURL(server.URL))
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

		//goroutine to periodically check if we need to stop a container
		go func() {
			ticker := time.NewTicker(time.Duration(routeStruct.HealthTimeout) * time.Millisecond)
			defer ticker.Stop()

			for range ticker.C {
				routeStruct.Lock.Lock()
				for i := len(routeStruct.Nodes) - 1; i >= 0; i-- {
					node := routeStruct.Nodes[i]
					node.Metrics.Lock.Lock()
					if routeStruct.Docker != nil && len(routeStruct.Nodes) > 1 && time.Since(node.Metrics.LastRequestTime).Milliseconds() > time.Duration(routeStruct.InactiveTimeout).Milliseconds() {
						delete(Balancer.NodeTable, node.ContainerID)
						routeStruct.RemoveNode(node)
					}
					node.Metrics.Lock.Unlock()
				}
				routeStruct.Lock.Unlock()
			}
		}()
	}

	return nil
}
