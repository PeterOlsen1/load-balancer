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
			logger.Info("Moving unhealthy node to inactive")
			p.mu.Lock()
			p.unsafeRemoveActive(n)
			p.unsafeAddInactive(n)
			p.mu.Unlock()
		}
	}

	for _, n := range p.Inactive {
		res, err := n.CheckHealth()
		if res == "healthy" && err == nil {
			logger.Info("Moving healthy node to active")
			p.mu.Lock()
			p.unsafeRemoveInactive(n)
			p.unsafeAddActive(n)
			p.mu.Unlock()
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

func (p *NodePool) UnpauseOne() error {
	if len(p.Inactive) == 0 {
		return fmt.Errorf("inactive pool empty")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

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

	p.mu.Lock()
	defer p.mu.Unlock()

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
	return len(p.Inactive)
}

func (p *NodePool) AddInactive(n *node.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Inactive = append(p.Inactive, n)
}

// This method does not lock the `p.mu` before performing
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

// This method does not lock the `p.mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeAddActive(n *node.Node) {
	p.Active = append(p.Active, n)
}

// This method does not lock the `p.mu` before performing
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

// This method does not lock the `p.mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeAddInactive(n *node.Node) {
	p.Inactive = append(p.Inactive, n)
}
