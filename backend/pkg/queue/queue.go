package queue

import (
	"errors"
	"load-balancer/pkg/types"
	"sync"
)

//add shared queue implementation here

var ConnectionQueue = &Queue{
	Notify: make(chan struct{}, 1),
}

type Queue struct {
	connections []*types.Connection
	lock        sync.Mutex
	Notify      chan struct{}
}

func (q *Queue) Dequeue() (*types.Connection, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.connections == nil {
		return nil, errors.New("cannot dequeue from uninitialized queue")
	}

	if len(q.connections) == 0 {
		return nil, errors.New("cannot dequeue from empty queue")
	}

	conn := q.connections[0]
	q.connections = q.connections[1:]
	return conn, nil
}

func (q *Queue) Enqueue(conn *types.Connection) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.connections = append(q.connections, conn)

	//i guess this is how go developers signal a channel without blocking
	select {
	case q.Notify <- struct{}{}:
		// Sent signal if buffer isn't full
	default:
		// Do nothing if buffer is full
	}
}
