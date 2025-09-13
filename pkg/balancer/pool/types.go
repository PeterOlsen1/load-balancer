package pool

import (
	"load-balancer/pkg/balancer/node"
	"sync"
)

type NodePool struct {
	Active   []*node.Node
	Inactive []*node.Node
	Heap     *NodeHeap
	mu       sync.Mutex
	isClosed bool
}

type NodeHeap struct {
	mu   sync.Mutex
	heap []*node.Node
}
