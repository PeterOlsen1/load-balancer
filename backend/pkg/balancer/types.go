package balancer

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
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
	Nodes  []*node.Node
	lock   sync.Mutex
	Routes []*Route
}

//make hash table between container ID and what route they live in for easy lookup?

type Route struct {
	config.RouteConfig
	lock  sync.Mutex
	Nodes []*node.Node
}
