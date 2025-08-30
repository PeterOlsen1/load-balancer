package pool

import (
	"crypto/sha256"
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"sync"
)

var roundRobinIndex = 0
var roundRobinIndexMu sync.Mutex

func (p *NodePool) RoundRobin() *node.Node {
	n := p.GetActiveSize()
	nodes := p.GetActive()

	if n == 0 {
		logger.Err("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		ws.EventEmitter.Error("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		return nil
	}

	roundRobinIndexMu.Lock()
	node := nodes[roundRobinIndex%n]
	roundRobinIndex++
	roundRobinIndex %= n
	roundRobinIndexMu.Unlock()

	loops := 0
	for node.Metrics.Health != "healthy" {
		node = nodes[roundRobinIndex%n]
		roundRobinIndex++
		roundRobinIndex %= n

		if loops > n {
			logger.Err("Could not find node to proxy", fmt.Errorf("found no healthy nodes"))
			return nil
		} else {
			loops++
		}
	}

	return node
}

func (p *NodePool) LeastConnections() *node.Node {
	var lowest *node.Node = nil
	for _, n := range p.GetActive() {
		if n.Queue.Len() < lowest.Queue.Len() && n.Metrics.Health != "healthy" {
			lowest = n
		}
	}
	return lowest
}

// this would be a little tougher, all docker containers are
// on my local machine, so should have same compute
func (p *NodePool) ComputeBased() *node.Node {
	return nil
}

func (p *NodePool) IpHash(ip string) *node.Node {
	nodes := p.GetActive()
	hash := sha256.Sum256([]byte(ip))
	hashInt := int(hash[0])
	idx := hashInt % len(nodes)
	node := nodes[idx]

	for node.Metrics.Health != "healthy" {
		hash := sha256.Sum256([]byte(ip))
		hashInt := int(hash[0])
		idx := hashInt % len(nodes)
		node = nodes[idx]
	}

	return node
}
