package route

import (
	"load-balancer/pkg/balancer/pool"
	"load-balancer/pkg/config"
	"load-balancer/pkg/types"
	"sync"
	"time"
)

type Route struct {
	config.RouteConfig
	NodePool  *pool.NodePool
	Queue     *RouteQueue
	LastScale time.Time
}

type RouteQueue struct {
	Lock       sync.Mutex          `json:"-"`
	Queue      []*types.Connection `json:"queue"`
	connSignal chan (struct{})
}
