package node

import (
	"os/exec"
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
	DockerInfo *DockerInfo `json:"docker_info"`
	Address    string      `json:"address"`
	Metrics    NodeMetrics `json:"metrics"`
}

type DockerInfo struct {
	Cmd *exec.Cmd
	Id  string `json:"id"`
}

type NodeMetrics struct {
	Lock         sync.Mutex
	Health       NodeHealth `json:"health"`
	ResponseTime float32    `json:"response_time"`
	Connections  uint32     `json:"connections"`
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
