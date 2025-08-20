package node

import (
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
}

type NodeMetrics struct {
	Lock         sync.Mutex `json:"-"`
	Health       string     `json:"health"`
	ResponseTime float32    `json:"response_time"`
	Connections  uint32     `json:"connections"`
}
