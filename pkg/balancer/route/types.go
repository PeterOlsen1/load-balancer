package route

import (
	"load-balancer/pkg/balancer/pool"
	"load-balancer/pkg/config"
	"load-balancer/pkg/types"
	"time"
)

type Route struct {
	config.RouteConfig
	NodePool  *pool.NodePool
	Queue     *RouteQueue
	LastScale time.Time
}

type RouteQueue struct {
	Queue chan *types.Connection
}
