package pool

import "load-balancer/pkg/balancer/node"

func InitPool() *NodePool {
	return &NodePool{
		Active:   make([]*node.Node, 0),
		Inactive: make([]*node.Node, 0),
	}
}
