package balancer

import "crypto/sha256"

var roundRobinIndex = 0

func (b *Balancer) RoundRobin() *Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	idx := roundRobinIndex % len(b.nodes)
	node := b.nodes[idx]
	roundRobinIndex++

	for node.Metrics.Health == Unhealthy {
		idx := roundRobinIndex % len(b.nodes)
		node = b.nodes[idx]
		roundRobinIndex++
	}

	return node
}

func (b *Balancer) LeastConnections() *Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	var lowest *Node = nil
	for _, n := range b.nodes {
		if n.Metrics.Connections < lowest.Metrics.Connections && n.Metrics.Health != Unhealthy {
			lowest = n
		}
	}
	return lowest
}

// this would be a little tougher, all docker containers are
// on my local machine, so should have same compute
func (b *Balancer) ComputeBased() *Node {
	return nil
}

func (b *Balancer) IPHash(ip string) *Node {
	hash := sha256.Sum256([]byte(ip))
	hashInt := int(hash[0])
	idx := hashInt % len(b.nodes)
	node := b.nodes[idx]

	for node.Metrics.Health == Unhealthy {
		hash := sha256.Sum256([]byte(ip))
		hashInt := int(hash[0])
		idx := hashInt % len(b.nodes)
		node = b.nodes[idx]
	}

	return node
}
