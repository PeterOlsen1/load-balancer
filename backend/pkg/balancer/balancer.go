package balancer

import (
	"load-balancer/pkg/queue"
	"time"
)

var PORT int = 3001
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
func (b *Balancer) InitBalancer(healthCheckPeriod int) {
	node, err := StartServer(3000)
	if err != nil {
		//error is already logged in the StartServer function
		return
	}

	//allow the server to start up before sending health request
	time.Sleep(2 * time.Second)
	b.AddNode(node)

	if healthCheckPeriod <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(healthCheckPeriod) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			b.lock.Lock()
			for _, n := range LoadBalancer.Nodes {
				n.CheckHealth()
			}
			b.lock.Unlock()
		}
	}()
}

func (b *Balancer) Lock() {
	b.lock.Lock()
}

func (b *Balancer) Unlock() {
	b.lock.Unlock()
}
