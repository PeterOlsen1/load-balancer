package node

import (
	"fmt"
	"load-balancer/pkg/types"
)

func NewNodeQueue(capacity int) *NodeQueue {
	return &NodeQueue{
		Queue: make([]*types.Connection, 0, capacity),
		Open:  true,
	}
}

func (q *NodeQueue) Enqueue(conn *types.Connection) error {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	if len(q.Queue) == cap(q.Queue) {
		return fmt.Errorf("queue is at capacity")
	}

	q.Queue = append(q.Queue, conn)
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

func (q *NodeQueue) CloseQueue() {
	q.Lock.Lock()
	q.Open = false
	q.Lock.Unlock()
}

func (q *NodeQueue) OpenQueue() {
	q.Lock.Lock()
	q.Open = true
	q.Lock.Unlock()
}
