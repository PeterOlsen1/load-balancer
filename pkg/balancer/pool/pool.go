package pool

import (
	"fmt"
	"load-balancer/pkg/balancer/node"
)

func (p *NodePool) CheckHealth() {
	fmt.Println(p.Active)
	fmt.Println(p.Inactive)

	for _, n := range p.Active {
		res, err := n.CheckHealth()
		if res != "healthy" || err != nil {
			fmt.Println("unhealthy -> inactive")
			p.mu.Lock()
			p.RemoveActive(n)
			p.AddInactive(n)
			p.mu.Unlock()
		}
	}

	for _, n := range p.Inactive {
		res, err := n.CheckHealth()
		if res == "healthy" && err == nil {
			fmt.Println("healthy -> active")
			p.mu.Lock()
			p.RemoveInactive(n)
			p.AddActive(n)
			p.mu.Unlock()
		}
	}
}

func (p *NodePool) GetAll() []*node.Node {
	out := make([]*node.Node, 0)

	for _, n := range p.Active {
		out = append(out, n)
	}
	for _, n := range p.Inactive {
		out = append(out, n)
	}

	return out
}

func (p *NodePool) GetActive() []*node.Node {
	return p.Active
}

func (p *NodePool) RemoveActive(n *node.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, node := range p.Active {
		if node == n {
			p.Active = append(p.Active[:i], p.Active[i+1:]...)
			break
		}
	}
}

func (p *NodePool) GetActiveSize() int {
	return len(p.Active)
}

func (p *NodePool) AddActive(n *node.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Active = append(p.Active, n)
}

func (p *NodePool) GetInactive() []*node.Node {
	return p.Inactive
}

func (p *NodePool) RemoveInactive(n *node.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, node := range p.Inactive {
		if node == n {
			p.Inactive = append(p.Inactive[:i], p.Inactive[i+1:]...)
			break
		}
	}
}

func (p *NodePool) GetInactiveSize() int {
	return len(p.Active)
}

func (p *NodePool) AddInactive(n *node.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Inactive = append(p.Inactive, n)
}
