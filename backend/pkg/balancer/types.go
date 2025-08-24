package balancer

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/balancer/route"
)

/*
The balancer is the main guy in this program.

It stores a list of all nodes that are to be accessed
from many different goroutines, hence the lock

Methods:
* AddNode(*Node) -> add a new node to the nodes list
* RemoveNode(*Node) -> remove a node from the nodes list
*/
type BalancerType struct {
	Routes    []*route.Route
	NodeTable map[string]*node.Node
}
