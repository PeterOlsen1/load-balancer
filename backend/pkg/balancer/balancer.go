package balancer

import (
	"load-balancer/pkg/logger"
	"load-balancer/pkg/queue"
	"os"
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

	f, err := os.OpenFile("./data/urls", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogErr("Failed to open url file", err)
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
