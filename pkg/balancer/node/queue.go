package node

import (
	"fmt"
	"load-balancer/pkg/batch"
	"load-balancer/pkg/types"
	"time"
)

func (n *Node) WatchQueue() {
	q := n.Queue
	batch := batch.InitBatch(100)
	batchTicker := time.NewTicker(time.Millisecond * 20)
	defer batchTicker.Stop()

	for {
		select {
		case <-batchTicker.C:
			for _, conn := range batch.Flush() {
				go n.processRequest(conn)
			}
		case conn := <-q.connSignal:
			if conn == nil {
				return
			}

			err := batch.Add(conn)
			if err != nil {
				for _, conn := range batch.Flush() {
					go n.processRequest(conn)
				}
			}
		case <-q.closeSignal:
			for conn := range q.connSignal {
				if conn == nil {
					break
				}
				go n.processRequest(conn)
			}
			for _, conn := range batch.Flush() {
				go n.processRequest(conn)
			}
			return
		}
	}
}

func InitNodeQueue(capacity int) *NodeQueue {
	return &NodeQueue{
		Queue:       make(chan *types.Connection, capacity),
		Open:        true,
		connSignal:  make(chan *types.Connection, capacity),
		closeSignal: make(chan struct{}),
	}
}

func (q *NodeQueue) Enqueue(conn *types.Connection) error {
	select {
	case q.connSignal <- conn:
		return nil
	default:
		return fmt.Errorf("queue is at capacity")
	}
}

func (q *NodeQueue) Dequeue() (*types.Connection, error) {
	select {
	case conn := <-q.connSignal:
		return conn, nil
	default:
		return nil, fmt.Errorf("queue is empty")
	}
}

func (n *Node) CloseQueue() {
	if !n.Queue.Open || n.Queue.closeSignal == nil {
		return
	}

	n.Queue.Open = false
	close(n.Queue.closeSignal)
	close(n.Queue.connSignal)
}

func (n *Node) OpenQueue() {
	n.Queue.Open = true
	n.Queue.connSignal = make(chan *types.Connection, cap(n.Queue.Queue))
	n.Queue.closeSignal = make(chan struct{})
	go n.WatchQueue()
}

func (q *NodeQueue) Len() int {
	return len(q.connSignal)
}
