package balancer

import (
	"crypto/sha256"
	"fmt"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/node"
)

var roundRobinIndex = 0

func (b *Balancer) RoundRobin() *node.Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.nodes) == 0 {
		logger.Err("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		return nil
	}

	idx := roundRobinIndex % len(b.nodes)
	node := b.nodes[idx]
	roundRobinIndex++

	for node.GetHealth() == "Unhealthy" {
		idx := roundRobinIndex % len(b.nodes)
		node = b.nodes[idx]
		roundRobinIndex++
	}

	return node
}

func (b *Balancer) LeastConnections() *node.Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	var lowest *node.Node = nil
	for _, n := range b.nodes {
		if n.Metrics.Connections < lowest.Metrics.Connections && n.GetHealth() == "Unhealthy" {
			lowest = n
		}
	}
	return lowest
}

// this would be a little tougher, all docker containers are
// on my local machine, so should have same compute
func (b *Balancer) ComputeBased() *node.Node {
	return nil
}

func (b *Balancer) IPHash(ip string) *node.Node {
	hash := sha256.Sum256([]byte(ip))
	hashInt := int(hash[0])
	idx := hashInt % len(b.nodes)
	node := b.nodes[idx]

	for node.GetHealth() == "Unhealthy" {
		hash := sha256.Sum256([]byte(ip))
		hashInt := int(hash[0])
		idx := hashInt % len(b.nodes)
		node = b.nodes[idx]
	}

	return node
}
