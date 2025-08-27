package route

import (
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"load-balancer/pkg/types"
	"sync"
)

type Route struct {
	config.RouteConfig
	Lock  sync.Mutex
	Nodes []*node.Node
	Queue *RouteQueue
}

type RouteQueue struct {
	Lock       sync.Mutex          `json:"-"`
	Queue      []*types.Connection `json:"queue"`
	connSignal chan (struct{})
}
