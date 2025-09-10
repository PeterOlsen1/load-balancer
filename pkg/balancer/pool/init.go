package pool

import (
	"container/heap"
	"load-balancer/pkg/balancer/node"
)

func InitPool() *NodePool {
	nodeHeap := &NodeHeap{}
	heap.Init(nodeHeap)

	return &NodePool{
		Active:   make([]*node.Node, 0),
		Inactive: make([]*node.Node, 0),
		Heap:     nodeHeap,
	}
}
