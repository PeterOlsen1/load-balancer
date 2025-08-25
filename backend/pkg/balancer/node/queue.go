package node

import (
	"fmt"
	"load-balancer/pkg/types"
)

func (n *Node) WatchQueue() {
	q := &n.Queue

	for {
		select {
		case <-q.signal:
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
	return &NodeQueue{
		Queue:  make([]*types.Connection, 0, capacity),
		Open:   true,
		signal: make(chan struct{}),
	}
}

func (q *NodeQueue) Enqueue(conn *types.Connection) error {
	q.Lock.Lock()

	if len(q.Queue) == cap(q.Queue) {
		q.Lock.Unlock()
		return fmt.Errorf("queue is at capacity")
	}

	q.Queue = append(q.Queue, conn)
	q.Lock.Unlock()

	q.signal <- struct{}{}
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
	close(n.Queue.signal)
	n.Queue.Lock.Unlock()
}

func (n *Node) OpenQueue() {
	n.Queue.Lock.Lock()
	n.Queue.Open = true
	n.Queue.signal = make(chan struct{})
	n.Queue.closeSignal = make(chan struct{})
	n.Queue.Lock.Unlock()
	go n.WatchQueue()
}
