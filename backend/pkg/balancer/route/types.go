package route

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"sync"
)

type Route struct {
	config.RouteConfig
	Lock  sync.Mutex
	Nodes []*node.Node
}
