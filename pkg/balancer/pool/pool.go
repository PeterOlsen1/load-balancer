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
	if p.isClosed {
		return
	}

	for _, n := range p.Active {
		go func(n *node.Node) {
			res, err := n.CheckHealth()
			if res != "healthy" || err != nil {
				logger.Info(fmt.Sprintf("Moving unhealthy node to inactive: %s", n.Address))
				p.mu.Lock()
				p.unsafeRemoveActive(n)
				p.unsafeAddInactive(n)
				p.mu.Unlock()
			}
		}(n)
	}

	for _, n := range p.Inactive {
		go func(n *node.Node) {
			res, err := n.CheckHealth()
			if res == "healthy" && err == nil {
				logger.Info(fmt.Sprintf("Moving healthy node to active: %s", n.Address))
				p.mu.Lock()
				p.unsafeRemoveInactive(n)
				p.unsafeAddActive(n)

				logger.PoolSize(len(p.Active), len(p.Inactive))
				p.mu.Unlock()
			}
		}(n)
	}

	// if p.GetActiveSize() < cfg.Pool.ActiveSize {
	// 	diff := cfg.Pool.ActiveSize - p.GetActiveSize()
	// 	logger.Info(fmt.Sprintf("Pool has fewer active nodes than config, unpausing %d", diff))

	// 	for range diff {
	// 		p.UnpauseOne()
	// 	}
	// }
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

func (p *NodePool) GetActiveSize() uint16 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return uint16(len(p.Active))
}

func (p *NodePool) AddActive(n *node.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Active = append(p.Active, n)
	p.Heap.Push(n)
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

			//remove from inactive, add to active
			p.Inactive = append(p.Inactive[:i], p.Inactive[i+1:]...)
			p.unsafeAddActive(n)
			if !n.Queue.IsOpen() {
				n.OpenQueue()
			}

			logger.Info(fmt.Sprintf("Unpaused one node: %s", n.Address))
			break
		}
	}

	logger.PoolSize(len(p.Active), len(p.Inactive))

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
	if n.Queue.IsOpen() {
		n.CloseQueue()
	}

	logger.Info(fmt.Sprintf("Paused one node: %s", n.Address))
	logger.PoolSize(len(p.Active), len(p.Inactive))

	return nil
}

func (p *NodePool) GetInactive() []*node.Node {
	return p.Inactive
}

func (p *NodePool) GetInactiveSize() uint16 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return uint16(len(p.Inactive))
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
	p.Heap.RemoveNode(n)
}

// This method does not lock the `p.mu` before performing
// its operation, and is therefore unsafe.
//
// Only use when the calling method acquires a lock
func (p *NodePool) unsafeAddActive(n *node.Node) {
	p.Active = append(p.Active, n)
	p.Heap.Add(n)
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

func (p *NodePool) Close() {
	p.isClosed = true

	// this is only called when shutting down, this is okay
	p.mu.Lock()
}
