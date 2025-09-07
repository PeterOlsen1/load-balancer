package pool

import (
	"load-balancer/pkg/balancer/node"
	"sync"
)

type NodePool struct {
	Active   []*node.Node
	Inactive []*node.Node
	Mu       sync.Mutex
}
