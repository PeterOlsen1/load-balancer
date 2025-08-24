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

func ConsumeMultiplePorts(numPorts int) []int {
	var ret []int
	portMutex.Lock()
	for range numPorts {
		ret = append(ret, port)
		port++
	}
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
			serverNode, err := routeStruct.Scale()
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

		//goroutine to periodically check health of containers
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
				routeStruct.lock.Lock()
				for _, n := range routeStruct.Nodes {
					go n.CheckHealth()
				}
				routeStruct.lock.Unlock()
			}
		}()

		//goroutine to periodically check if we need to stop a container
		go func() {
			ticker := time.NewTicker(time.Duration(route.HealthTimeout) * time.Millisecond)
			defer ticker.Stop()

			for range ticker.C {
				routeStruct.lock.Lock()
				for i := len(routeStruct.Nodes) - 1; i >= 0; i-- {
					node := routeStruct.Nodes[i]
					node.Metrics.Lock.Lock()
					if routeStruct.Docker != nil && len(routeStruct.Nodes) > 1 && time.Since(node.Metrics.LastRequestTime).Milliseconds() > time.Duration(routeStruct.InactiveTimeout).Milliseconds() {
						routeStruct.RemoveNode(node)
					}
					node.Metrics.Lock.Unlock()
				}
				routeStruct.lock.Unlock()
			}
		}()
	}

	return nil
}
