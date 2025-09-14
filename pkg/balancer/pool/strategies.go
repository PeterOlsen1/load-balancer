package pool

import (
	"crypto/sha256"
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"sync"
)

var roundRobinIndex uint16 = 0
var roundRobinIndexMu sync.Mutex

func (p *NodePool) RoundRobin() *node.Node {
	n := p.GetActiveSize()

	if n == 0 {
		logger.Err("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		ws.EventEmitter.Error("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		return nil
	}

	nodes := p.GetActive()
	n = uint16(len(nodes))
	if n == 0 {
		return nil
	}
	roundRobinIndexMu.Lock()
	node := nodes[roundRobinIndex%n]
	roundRobinIndex++
	roundRobinIndex %= n
	roundRobinIndexMu.Unlock()

	var loops uint16 = 0
	// keep finding nodes until we find one that has space and is healthy
	for node.Metrics.Health != "healthy" || !node.Queue.HasSpace() {
		roundRobinIndexMu.Lock()
		node = nodes[roundRobinIndex%n]
		roundRobinIndex++
		roundRobinIndex %= n
		roundRobinIndexMu.Unlock()

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
	lowest, err := p.Heap.RemoveMin()
	if err != nil {
		logger.Err("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
		ws.EventEmitter.Error("Could not find node to proxy", fmt.Errorf("nodes length is 0"))
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
