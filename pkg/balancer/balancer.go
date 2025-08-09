package balancer

import (
	"fmt"
	"load-balancer/pkg/queue"
	"load-balancer/pkg/types"
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

			go handleQueuePop(conn)
		}
	}
}

func (b *Balancer) InitBalancer() {
	node, err := StartServer(3000)
	if err != nil {
		return
	}

	//allow the server to start up before sending health request
	time.Sleep(2 * time.Second)
	b.AddNode(node)
}

func handleQueuePop(conn *types.Connection) {
	fmt.Println(conn.Request.URL.Path)
}
