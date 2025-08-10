package balancer

import (
	"fmt"
	"load-balancer/pkg/logger"
	"net/http"
	"os/exec"
	"strings"
)

// Helper method to start an internal server,
//
// In a real environment, this would not be necessary,
// and the user would just call the Balancer.AddNode method
func StartServer(port int) (*Node, error) {
	path := "./server/run.sh" //assuming you run from root of project

	cmd := exec.Command("bash", path, fmt.Sprintf("%d", port))

	output, err := cmd.Output()
	if err != nil {
		go logger.LogErr("Creating container", err)
		return nil, err
	}
	containerID := strings.TrimSpace(string(output))
	if containerID == "" {
		err := fmt.Errorf("empty container ID received")
		go logger.LogErr("Creating container", err)
		return nil, err
	}
	go logger.LogContainerStart(containerID)

	node := Node{
		DockerInfo: &DockerInfo{
			Cmd: cmd,
			id:  containerID,
		},
		Address: fmt.Sprintf("http://localhost:%d", port),
	}

	go logger.Log(fmt.Sprintf("Started server @ http://localhost: %d", port))
	return &node, nil
}

func (b *Balancer) AddNode(node *Node) {
	b.lock.Lock()
	defer b.lock.Unlock()

	go node.CheckHealth()
	b.nodes = append(b.nodes, node)
}

// add response time metric
func (node *Node) CheckHealth() error {
	address := node.Address
	resp, err := http.Get(fmt.Sprintf("%s/health", address))
	if err != nil {
		go logger.LogErr("Fetching node health", err)
		return err
	}

	health := Healthy
	if resp.StatusCode != http.StatusOK {
		health = Unhealthy
		go logger.LogStatusCheck("Unhealthy", node.Address)
	} else {
		go logger.LogStatusCheck("Healthy", node.Address)
	}
	node.Metrics.Lock.Lock()
	defer node.Metrics.Lock.Unlock()
	node.Metrics.Health = health

	return nil
}

func (b *Balancer) RemoveNode(node *Node) error {
	node.StopServer()

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
	go logger.Log("cleaning up nodes")

	for _, n := range b.nodes {
		n.StopServer()
	}

	var empty []*Node
	b.nodes = empty
	return nil
}
