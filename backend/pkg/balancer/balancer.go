package balancer

import (
	"load-balancer/pkg/queue"
	"time"
)

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

func (b *Balancer) InitBalancer() {
	node, err := StartServer(3000)
	if err != nil {
		//error is already logged in the StartServer function
		return
	}

	//allow the server to start up before sending health request
	time.Sleep(2 * time.Second)
	b.AddNode(node)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			b.lock.Lock()
			for _, n := range LoadBalancer.nodes {
				n.CheckHealth()
			}
			b.lock.Unlock()
		}
	}()
}
