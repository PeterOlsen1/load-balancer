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
				q.workerPool.Event(conn)
			}
		case conn := <-q.connChan:
			if conn == nil {
				return
			}

			err := batch.Add(conn)
			if err != nil {
				for _, conn := range batch.Flush() {
					q.workerPool.Event(conn)
				}
			}
		case <-q.closeSignal:
			for conn := range q.connChan {
				if conn == nil {
					break
				}
				go n.processRequest(conn)
			}
			for _, conn := range batch.Flush() {
				go n.processRequest(conn)
			}

			q.workerPool.Close()
			return
		}
	}
}

func (q *NodeQueue) Enqueue(conn *types.Connection) error {
	select {
	case q.connChan <- conn:
		return nil
	default:
		return fmt.Errorf("queue is at capacity")
	}
}

func (q *NodeQueue) Dequeue() (*types.Connection, error) {
	select {
	case conn := <-q.connChan:
		return conn, nil
	default:
		return nil, fmt.Errorf("queue is empty")
	}
}

func (n *Node) CloseQueue() {
	if !n.Queue.open || n.Queue.closeSignal == nil {
		return
	}

	n.Queue.open = false
	close(n.Queue.closeSignal)
	close(n.Queue.connChan)
}

func (n *Node) OpenQueue() {
	n.Queue.open = true
	n.Queue.connChan = make(chan *types.Connection, cap(n.Queue.queue))
	n.Queue.closeSignal = make(chan struct{})
	go n.WatchQueue()
}

func (q *NodeQueue) Len() int {
	return len(q.connChan)
}

func (q *NodeQueue) HasSpace() bool {
	return len(q.connChan) < cap(q.connChan)
}

func (q *NodeQueue) IsOpen() bool {
	return q.open
}
