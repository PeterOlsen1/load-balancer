package balancer

import (
	"os/exec"
	"sync"
)

/*
The balancer is the main guy in this program.

It stores a list of all nodes that are to be accessed
from many different goroutines, hence the lock

Methods:
* AddNode(*Node) -> add a new node to the nodes list
* RemoveNode(*Node) -> remove a node from the nodes list
*/
type Balancer struct {
	//nodes + node health?
	nodes []*Node
	lock  sync.Mutex
}

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
	DockerInfo *DockerInfo
	Address    string
	Metrics    NodeMetrics
}

type DockerInfo struct {
	Cmd *exec.Cmd
	id  string
}

type NodeMetrics struct {
	Lock         sync.Mutex
	Health       NodeHealth
	ResponseTime float32
	Connections  uint32
}

/*
Enum to keep track of node health,
taken from the `/health` route of the servers

Unknown: uninitialized
Unhealthy: bad status code was returned
Healthy: 2** status code was returned
*/
type NodeHealth int

const (
	Unknown NodeHealth = iota
	Unhealthy
	Healthy
)
