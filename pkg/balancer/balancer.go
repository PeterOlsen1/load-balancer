package balancer

import (
	"fmt"
	"load-balancer/pkg/queue"
	"load-balancer/pkg/types"
)

var balancer = Balancer{}

type Balancer struct {
	//nodes + node health?
}

func WatchQueue() {
	for {
		<- queue.ConnectionQueue.Notify
		for {
			conn, err := queue.ConnectionQueue.Dequeue()

			if err != nil {
				break
			}

			go handleQueuePop(conn)
		}
	}
}

func handleQueuePop(conn *types.Connection) {
	fmt.Println(conn.Request.URL.Path)
}