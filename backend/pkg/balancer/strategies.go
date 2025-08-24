package balancer

import (
	"crypto/sha256"
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
)

var roundRobinIndex = 0

func (r *Route) GetProxyNode(ip string) *node.Node {
	switch r.Strategy {
	case "round-robin":
		return r.roundRobin()
	case "least-connections":
		return r.leastConnections()
	case "compute-based":
		return r.computeBased()
	case "ip-hash":
		return r.ipHash(ip)
	default:
		return nil
	}
}

func (r *Route) roundRobin() *node.Node {
	r.lock.Lock()
	if len(r.Nodes) == 0 {
		logger.Err("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		ws.EventEmitter.Error("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		return nil
	}

	idx := roundRobinIndex % len(r.Nodes)
	node := r.Nodes[idx]
	roundRobinIndex++
	r.lock.Unlock()

	n := len(r.Nodes)
	loops := 0
	for node.Metrics.Health != "healthy" {
		idx := roundRobinIndex % len(r.Nodes)
		node = r.Nodes[idx]
		roundRobinIndex++
		
		if loops > n {
			return nil
		} else {
			loops++
		}
	}

	return node
}

func (r *Route) leastConnections() *node.Node {
	r.lock.Lock()
	defer r.lock.Unlock()

	var lowest *node.Node = nil
	for _, n := range r.Nodes {
		if n.Metrics.Connections < lowest.Metrics.Connections && n.Metrics.Health != "healthy" {
			lowest = n
		}
	}
	return lowest
}

// this would be a little tougher, all docker containers are
// on my local machine, so should have same compute
func (r *Route) computeBased() *node.Node {
	return nil
}

func (r *Route) ipHash(ip string) *node.Node {
	hash := sha256.Sum256([]byte(ip))
	hashInt := int(hash[0])
	idx := hashInt % len(r.Nodes)
	node := r.Nodes[idx]

	for node.Metrics.Health != "healthy" {
		hash := sha256.Sum256([]byte(ip))
		hashInt := int(hash[0])
		idx := hashInt % len(r.Nodes)
		node = r.Nodes[idx]
	}

	return node
}
