package balancer

import (
	"fmt"
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
* CheckNode(*Node) -> send a request to the given node, and update its metrics
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

// Stops the server associated with any given node
// through the docker stop command.
//
// If this node has no server, instantly return nil
func (node *Node) StopServer() error {
	if node.DockerInfo == nil {
		return nil
	}

	cmd := exec.Command("docker", "stop", node.DockerInfo.id)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error docker stop: ", err)
		return err
	}

	fmt.Println("Stopped container with ID: ", node.DockerInfo.id)
	return nil
}

type NodeMetrics struct {
	Lock         sync.Mutex
	Health       NodeHealth
	ResponseTime float64
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
