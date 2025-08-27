package route

import (
	"load-balancer/pkg/balancer/pool"
	"load-balancer/pkg/config"
	"load-balancer/pkg/types"
	"sync"
)

type Route struct {
	config.RouteConfig
	Lock     sync.Mutex
	NodePool *pool.NodePool
	Queue    *RouteQueue
}

type RouteQueue struct {
	Lock       sync.Mutex          `json:"-"`
	Queue      []*types.Connection `json:"queue"`
	connSignal chan (struct{})
}
