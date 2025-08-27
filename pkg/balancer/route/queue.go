package route

import (
	"fmt"
	"load-balancer/pkg/errors"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/types"
	"time"
)

func (r *Route) WatchQueue() {
	q := r.Queue

	for range q.connSignal {
		go func() {
			conn, err := q.Dequeue()

			if err != nil {
				logger.Err("Failed to dequeue from route queue", err)
				errors.Send500(conn, "Failed to dequeue from route queue")
				return
			}

			node := r.GetProxyNode(conn.Request.RemoteAddr)
			if node == nil {
				logger.Err("Failed to find node for proxy", fmt.Errorf("failed to find node for proxy"))
				errors.Send500(conn, "Failed to find node for proxy")
				return
			}

			err = node.Queue.Enqueue(conn)
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

				q.EnqueueFront(conn)
				return
			}

			//do some proprietary health checks of the node down here
			// add new node if we are above x connections
			// if we have one connection (slow) and more than one node, remove it
			// ^ could be improved upon,
			if !node.Metrics.CreatedNewNode && node.Queue.Len() > r.Docker.RequestScaleThreshold {
				node.Metrics.CreatedNewNode = true
				go func() {
					_, err := r.Scale()
					if err != nil {
						errors.Send500(conn, "Failed starting server on connection threshhold")
						return
					}

					// b.NodeTable[node.ContainerID] = node
				}()
			}

			// if we are below 70% of connection threshold, it is okay to make a new node
			// if len(node.Queue.Queue) < int(float64(routeObject.Docker.RequestScaleThreshold)*0.7) {
			// node.Metrics.CreatedNewNode = false
			// }

			node.Metrics.Lock.Lock()
			node.Metrics.LastRequestTime = time.Now()
			node.Metrics.Lock.Unlock()
		}()
	}
}

func InitRouteQueue() *RouteQueue {
	var q []*types.Connection
	q = make([]*types.Connection, 0)

	return &RouteQueue{
		Queue:      q,
		connSignal: make(chan struct{}),
	}
}

func (q *RouteQueue) Enqueue(conn *types.Connection) {
	q.Lock.Lock()
	q.Queue = append(q.Queue, conn)
	q.Lock.Unlock()

	q.connSignal <- struct{}{}
}

func (q *RouteQueue) Dequeue() (*types.Connection, error) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	if len(q.Queue) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	conn := q.Queue[0]
	q.Queue = q.Queue[1:]
	return conn, nil
}

func (q *RouteQueue) EnqueueFront(conn *types.Connection) {
	q.Lock.Lock()
	q.Queue = append([]*types.Connection{conn}, q.Queue...)
	q.Lock.Unlock()

	q.connSignal <- struct{}{}
}
