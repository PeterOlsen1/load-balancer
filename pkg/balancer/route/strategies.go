package route

import (
	"load-balancer/pkg/balancer/node"
)

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
