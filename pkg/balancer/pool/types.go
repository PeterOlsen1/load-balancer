package pool

import (
	"load-balancer/pkg/balancer/node"
	"sync"
)

type NodePool struct {
	Active   []*node.Node
	Inactive []*node.Node
	Heap     *NodeHeap
	Mu       sync.Mutex
}

type NodeHeap struct {
	mu   sync.Mutex
	heap []*node.Node
}
