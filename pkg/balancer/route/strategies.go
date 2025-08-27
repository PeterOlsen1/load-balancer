package route

import (
	"load-balancer/pkg/balancer/node"
	"sync"
)

// this WILL be bad once we get more routes. cross that bridge when it comes.
// add round robin object to each route?
var roundRobinIndex = 0
var roundRobinIndexMu sync.Mutex

func (r *Route) GetProxyNode(ip string) *node.Node {
	switch r.Strategy {
	case "round-robin":
		return r.NodePool.RoundRobin()
	case "least-connections":
		return r.NodePool.LeastConnections()
	case "compute-based":
		return r.NodePool.ComputeBased()
	case "ip-hash":
		return r.NodePool.IpHash(ip)
	default:
		return nil
	}
}
