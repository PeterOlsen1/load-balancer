package pool

import (
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
)

// Move all unhealthy nodes in active to inactive
// Move all healthy nodes in inactive in active
//
// Paused nodes in inactive will not be moved
func (p *NodePool) CheckHealth(cfg config.RouteConfig) {
	for _, n := range p.Active {
		res, err := n.CheckHealth()
		if res != "healthy" || err != nil {
			logger.Info(fmt.Sprintf("Moving unhealthy node to inactive: %s", n.Address))
			p.Mu.Lock()
			p.unsafeRemoveActive(n)
			p.unsafeAddInactive(n)
			p.Mu.Unlock()
		}
	}

	for _, n := range p.Inactive {
		res, err := n.CheckHealth()
		if res == "healthy" && err == nil {
			logger.Info(fmt.Sprintf("Moving healthy node to active: %s", n.Address))
			p.Mu.Lock()
			p.unsafeRemoveInactive(n)
			p.unsafeAddActive(n)
			p.Mu.Unlock()
		}
	}

	if p.GetActiveSize() < cfg.Pool.ActiveSize {
		diff := cfg.Pool.ActiveSize - p.GetActiveSize()
		logger.Info(fmt.Sprintf("Pool has fewer active nodes than config, unpausing %d", diff))

		for range diff {
			p.UnpauseOne()
		}
	}
}

func (p *NodePool) GetAll() []*node.Node {
	out := make([]*node.Node, 0)

	out = append(out, p.Active...)
	out = append(out, p.Inactive...)
	return out
}

func (p *NodePool) GetActive() []*node.Node {
	return p.Active
}

func (p *NodePool) RemoveActive(n *node.Node) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	for i, node := range p.Active {
		if node == n {
			p.Active = append(p.Active[:i], p.Active[i+1:]...)
			break
		}
	}
}

func (p *NodePool) GetActiveSize() int {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	return len(p.Active)
}

func (p *NodePool) AddActive(n *node.Node) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Active = append(p.Active, n)
}

func (p *NodePool) UnpauseOne() error {
	if len(p.Inactive) == 0 {
		return fmt.Errorf("inactive pool empty")
	}

	p.Mu.Lock()
	defer p.Mu.Unlock()

	//loop through inactive nodes, health check, activate the first good one
	for i, n := range p.Inactive {
		if n.Metrics.Health != "unhealthy" {
			n.Metrics.Health = "unknown" //set to unknown so health check doesn't insta-return from pause
			health, err := n.CheckHealth()
			if err != nil || health != "healthy" {
				continue
			}

			p.Inactive = append(p.Inactive[:i], p.Inactive[i+1:]...)
			p.unsafeAddActive(n)
			if !n.Queue.Open {
				n.OpenQueue()
			}

			logger.Info(fmt.Sprintf("Unpaused one node: %s", n.Address))
			return nil
		}
	}

	return nil
}

func (p *NodePool) PauseOne() error {
	if len(p.Active) == 0 {
		return fmt.Errorf("active pool empty")
	}

	p.Mu.Lock()
	defer p.Mu.Unlock()

	//loop through inactive nodes, health check, activate the first good one
	n := p.Active[0]
	n.Metrics.Health = "paused"

	p.Active = p.Active[1:]
	p.unsafeAddInactive(n)
	if n.Queue.Open {
		n.CloseQueue()
	}

	logger.Info(fmt.Sprintf("Paused one node: %s", n.Address))

	return nil
}

func (p *NodePool) GetInactive() []*node.Node {
	return p.Inactive
}

func (p *NodePool) RemoveInactive(n *node.Node) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	for i, node := range p.Inactive {
		if node == n {
			p.Inactive = append(p.Inactive[:i], p.Inactive[i+1:]...)
			break
		}
	}
}

func (p *NodePool) GetInactiveSize() int {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	return len(p.Inactive)
}

func (p *NodePool) AddInactive(n *node.Node) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Inactive = append(p.Inactive, n)
}

// This method does not lock the `p.Mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeRemoveActive(n *node.Node) {
	for i, node := range p.Active {
		if node == n {
			p.Active = append(p.Active[:i], p.Active[i+1:]...)
			break
		}
	}
}

// This method does not lock the `p.Mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeAddActive(n *node.Node) {
	p.Active = append(p.Active, n)
}

// This method does not lock the `p.Mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeRemoveInactive(n *node.Node) {
	for i, node := range p.Inactive {
		if node == n {
			p.Inactive = append(p.Inactive[:i], p.Inactive[i+1:]...)
			break
		}
	}
}

// This method does not lock the `p.Mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeAddInactive(n *node.Node) {
	p.Inactive = append(p.Inactive, n)
}
