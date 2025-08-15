package balancer

import (
	"crypto/sha256"
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
)

var roundRobinIndex = 0

func (b *Balancer) RoundRobin() *node.Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.Nodes) == 0 {
		go logger.Err("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		go ws.EventEmitter.Error("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		return nil
	}

	idx := roundRobinIndex % len(b.Nodes)
	node := b.Nodes[idx]
	roundRobinIndex++

	for node.Metrics.Health != "healthy" {
		idx := roundRobinIndex % len(b.Nodes)
		node = b.Nodes[idx]
		roundRobinIndex++
	}

	return node
}

func (b *Balancer) LeastConnections() *node.Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	var lowest *node.Node = nil
	for _, n := range b.Nodes {
		if n.Metrics.Connections < lowest.Metrics.Connections && n.Metrics.Health != "healthy" {
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
	idx := hashInt % len(b.Nodes)
	node := b.Nodes[idx]

	for node.Metrics.Health != "healthy" {
		hash := sha256.Sum256([]byte(ip))
		hashInt := int(hash[0])
		idx := hashInt % len(b.Nodes)
		node = b.Nodes[idx]
	}

	return node
}
