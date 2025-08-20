package balancer

import (
	"fmt"
	"load-balancer/pkg/config"
	"load-balancer/pkg/queue"
	"time"
)

var PORT int = 3000
var LoadBalancer = Balancer{}

func WatchQueue() {
	for {
		<-queue.ConnectionQueue.Notify
		for {
			conn, err := queue.ConnectionQueue.Dequeue()

			if err != nil {
				break
			}

			LoadBalancer.ProxyRequest(conn)
		}
	}
}

// Pass in num <= 0 for no health checks
func (b *Balancer) InitBalancer(healthCheckPeriod int) error {
	for _, route := range config.Config.Routes {
		routeStruct := Route{
			RouteConfig: route,
		}

		b.Routes = append(b.Routes, &routeStruct)
		if route.Docker == nil {
			continue
		}

		node, err := StartServer(PORT, route.Docker)
		if err != nil {
			return err
		}

		routeStruct.lock.Lock()
		routeStruct.Nodes = append(routeStruct.Nodes, node)
		routeStruct.lock.Unlock()
	}

	//allow the server to start up before sending health request
	time.Sleep(1500 * time.Millisecond)

	if healthCheckPeriod <= 0 {
		return fmt.Errorf("health check period is negative")
	}

	go func() {
		ticker := time.NewTicker(time.Duration(healthCheckPeriod) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			b.lock.Lock()
			for _, n := range LoadBalancer.Nodes {
				n.CheckHealth()
			}
			b.lock.Unlock()
		}
	}()

	return nil
}