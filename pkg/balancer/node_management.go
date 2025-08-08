package balancer

import (
	"fmt"
	"net/http"
	"os/exec"
)

// Helper method to start an internal server,
//
// In a real environment, this would not be necessary,
// and the user would just call the Balancer.AddNode method
func StartServer(port int) (*Node, error) {
	path := "../../server/run.sh"

	cmd := exec.Command("bash", path, fmt.Sprintf("%d", port))
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start docker container: ", err)
		return nil, err
	}

	node := Node{
		Cmd:     cmd,
		Address: fmt.Sprintf("localhost:%d", port),
	}
	return &node, nil
}

func (b *Balancer) AddNode(node *Node) {
	b.lock.Lock()
	defer b.lock.Unlock()

	go b.CheckNode(node)
	b.nodes = append(b.nodes, node)
}

// add response time metric
func (b *Balancer) CheckNode(node *Node) error {
	address := node.Address
	resp, err := http.Get(fmt.Sprintf("%s/health", address))
	if err != nil {
		fmt.Println("Error fetching node health: ", err)
		return err
	}

	health := Healthy
	if resp.StatusCode != http.StatusOK {
		health = Unhealthy
	}
	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()
	node.Metrics.Health = health

	return nil
}

func (b *Balancer) RemoveNode(node *Node) error {
	if node.Cmd != nil {
		err := node.Cmd.Process.Kill()

		if err != nil {
			fmt.Println("Error killing process: ", err)
			return err
		}
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	var filtered []*Node
	for _, n := range b.nodes {
		if n != node {
			filtered = append(filtered, n)
		}
	}
	b.nodes = filtered

	return nil
}

func (b *Balancer) CleanupNodes() error {
	for _, n := range b.nodes {
		if n.Cmd != nil {
			err := n.Cmd.Process.Kill()

			if err != nil {
				fmt.Println("Error in cleanup: ", err)
				return err
			}
		}
	}

	var empty []*Node
	b.nodes = empty
	return nil
}
