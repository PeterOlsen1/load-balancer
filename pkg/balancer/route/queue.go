package route

import (
	"fmt"
	"load-balancer/pkg/errors"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"time"
)

func (r *Route) WatchQueue() {
	for conn := range r.Queue.Queue {
		go func(conn *types.Connection) {
			node := r.GetProxyNode(conn.Request.RemoteAddr)
			if node == nil {
				logger.Err("Failed to find node for proxy", fmt.Errorf("failed to find node for proxy"))
				errors.Send500(conn, "Failed to find node for proxy")
				return
			}

			err := node.Queue.Enqueue(conn)
			if err != nil {
				logger.Err("Node refused connection, retrying", err)
				conn.RetryCount++
				if conn.RetryCount > 3 {
					fmt.Println("retry limit exceeded")
					errors.Send500(conn, "Exceeded retry limit")
					return
				}

				//add a delay so that the same request isn't processed over and over again
				//set the duration to the response time of the last health check
				time.Sleep(time.Duration(node.Metrics.ResponseTime) * time.Millisecond)

				r.Queue.Queue <- conn
				return
			}

			r.NodePool.Heap.Add(node)
			load := r.CalculateLoad()
			if load > 50 {
				fmt.Println(load)
			}
			if load > 70 {
				r.Scale(r.RouteConfig)
			}
		}(conn)
	}
}

func InitRouteQueue(queueSize uint) *RouteQueue {
	return &RouteQueue{
		Queue: make(chan *types.Connection, queueSize),
	}
}

func (q *RouteQueue) Enqueue(conn *types.Connection) {
	q.Queue <- conn
}

func (q *RouteQueue) Dequeue() (*types.Connection, error) {
	conn := <-q.Queue
	return conn, nil
}

func (q *RouteQueue) Len() int {
	return len(q.Queue)
}
