package balancer

import (
	"load-balancer/pkg/types"
)

func (b *Balancer) ProxyRequest(conn *types.Connection) {
	node := b.RoundRobin()

	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()
	node.Metrics.Connections++

}
