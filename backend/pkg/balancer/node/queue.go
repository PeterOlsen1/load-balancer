package node

import (
	"fmt"
	"load-balancer/pkg/types"
)

func (n *Node) WatchQueue() {
	q := n.Queue

	for {
		select {
		case <-q.connSignal:
			conn, err := q.Dequeue()
			if err != nil {
				continue
			}

			n.processRequest(conn)
		case <-q.closeSignal:
			fmt.Println("closing queue, stop watching")
			for len(q.Queue) > 0 {
				conn, err := q.Dequeue()
				if err != nil {
					continue
				}

				go n.processRequest(conn)
			}
			return
		}
	}
}

func InitNodeQueue(capacity int) *NodeQueue {
	var q []*types.Connection
	if capacity > 0 {
		q = make([]*types.Connection, 0, capacity)
	} else {
		q = make([]*types.Connection, 0)
	}

	return &NodeQueue{
		Queue:       q,
		Open:        true,
		connSignal:  make(chan struct{}),
		closeSignal: make(chan struct{}),
	}
}

func (q *NodeQueue) Enqueue(conn *types.Connection) error {
	q.Lock.Lock()

	if len(q.Queue) >= cap(q.Queue) {
		q.Lock.Unlock()
		return fmt.Errorf("queue is at capacity")
	}

	q.Queue = append(q.Queue, conn)
	q.Lock.Unlock()

	q.connSignal <- struct{}{}
	return nil
}

func (q *NodeQueue) Dequeue() (*types.Connection, error) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	if len(q.Queue) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	conn := q.Queue[0]
	q.Queue = q.Queue[1:]
	return conn, nil
}

func (q *NodeQueue) TakeFromBack(numEntries int) ([]*types.Connection, error) {
	if numEntries == 0 {
		return nil, nil
	}

	q.Lock.Lock()
	defer q.Lock.Unlock()

	if len(q.Queue) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	if numEntries > len(q.Queue) {
		numEntries = len(q.Queue)
	}

	conns := make([]*types.Connection, numEntries)
	copy(conns, q.Queue[len(q.Queue)-numEntries:])

	q.Queue = q.Queue[:len(q.Queue)-numEntries]
	return conns, nil
}

func (n *Node) CloseQueue() {
	n.Queue.Lock.Lock()
	n.Queue.Open = false
	n.Queue.closeSignal <- struct{}{}

	close(n.Queue.closeSignal)
	close(n.Queue.connSignal)
	n.Queue.Lock.Unlock()
}

func (n *Node) OpenQueue() {
	n.Queue.Lock.Lock()
	n.Queue.Open = true
	n.Queue.connSignal = make(chan struct{})
	n.Queue.closeSignal = make(chan struct{})
	n.Queue.Lock.Unlock()
	go n.WatchQueue()
}

func (q *NodeQueue) Len() int {
	return len(q.Queue)
}
