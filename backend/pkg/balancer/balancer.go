package balancer

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"load-balancer/pkg/queue"
	"sync"
	"time"
)

var port int = 3000
var portMutex sync.Mutex

func ConsumePort() int {
	var ret int
	portMutex.Lock()
	ret = port
	port++
	portMutex.Unlock()
	return ret
}

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
	for _, route := range config.Config.Routes {
		routeStruct := Route{
			RouteConfig: route,
		}

		b.Routes = append(b.Routes, &routeStruct)
		if route.Docker != nil && len(route.Servers) == 0 {
			serverNode, err := StartServer(route.Docker)
			if err != nil {
				return err
			}

			routeStruct.lock.Lock()
			routeStruct.Nodes = append(routeStruct.Nodes, serverNode)
			routeStruct.lock.Unlock()
			b.NodeTable[serverNode.ContainerID] = serverNode
		}

		for _, server := range route.Servers {
			routeStruct.Nodes = append(routeStruct.Nodes, node.FromURL(server.URL))
		}

		go func() {
			if route.HealthTimeout <= 0 {
				//skip health check if timeout is not set
				return
			}

			//allow the server to start up before sending health request
			time.Sleep(1500 * time.Millisecond)

			ticker := time.NewTicker(time.Duration(route.HealthTimeout) * time.Millisecond)
			defer ticker.Stop()

			for range ticker.C {
				for _, r := range b.Routes {
					r.lock.Lock()
					for _, n := range r.Nodes {
						n.CheckHealth()
					}
					r.lock.Unlock()
				}
			}
		}()
	}

	return nil
}
