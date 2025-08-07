package queue

import (
	"errors"
	"load-balancer/pkg/types"
	"sync"
)

//add shared queue implementation here

var ConnectionQueue = &Queue{}

type Queue struct {
	connections []*types.Connection
	lock sync.Mutex
}

func (q *Queue) Pop() (*types.Connection, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.connections == nil {
		return nil, errors.New("cannot pop from uninitialized queue")
	}

	if len(q.connections) == 0 {
		return nil, errors.New("cannot pop from empty queue")
	}

	conn := q.connections[0]
	q.connections = q.connections[1:]
	return conn, nil
}

func (q *Queue) Push(conn *types.Connection) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.connections = append(q.connections, conn)
}