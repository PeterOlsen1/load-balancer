package node

import (
	"load-balancer/pkg/types"
	"load-balancer/pkg/workerpool"
	"sync"
)

/*
This struct defines the node type that all server
nodes will be stored as. The `Address` field defines
where all requests will be sent, and the given
server must contain `/health` route that returns
the health of the given node

Method:
* CheckHealth() -> send a request to the given node, and update its metrics
*/
type Node struct {
	ContainerID string      `json:"id"`
	Address     string      `json:"address"`
	Metrics     NodeMetrics `json:"metrics"`
	Queue       *NodeQueue  `json:"queue"`
	Weight      uint32      `json:"weight"`
}

// Explaining some fields
// * Response time is in ms
// * Health is an enum that can be any of the following values:
// - healthy
// - paused
// - unhealthy
// - unknown
// * CreatedNewNode is a flag where if we are over max connections, a node has already been created
type NodeMetrics struct {
	mu             sync.Mutex
	Health         string  `json:"health"`
	ResponseTime   float32 `json:"response_time"`
	Connections    uint32  `json:"connections"`
	CreatedNewNode bool    `json:"-"`
	// LastRequestTime time.Time  `json:"last_request"`
}

// A batch-processed request queue that implements a worker pool for processing
//
// This queue is where all connections are first sent before going directly to a node.
// This is done so that we can control the flow of requests,
// requeue to different nodes upon failure of this one,
// and easily calculate load level.
type NodeQueue struct {
	queue       chan *types.Connection // Channel-based queue
	open        bool                   // Indicates if the queue is open
	connChan    chan *types.Connection // Signal channel for new connections
	closeSignal chan struct{}          // Signal channel for closing the queue
	workerPool  *workerpool.WorkerPool[*types.Connection]
}
