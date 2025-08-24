package docker

import (
	"context"
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"
	"os/exec"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func StartServer(externalPort int, dockerInfo *config.DockerConfig) (*node.Node, error) {
	path := "./server/run.sh" //assuming you run from root of project

	cmd := exec.Command("bash", path, dockerInfo.Image, fmt.Sprintf("%d", externalPort), fmt.Sprintf("%d", dockerInfo.InternalPort))

	output, err := cmd.Output()
	if err != nil {
		logger.Err("Creating container", err)
		ws.EventEmitter.Error("Creating container", err)
		return nil, err
	}
	containerID := strings.TrimSpace(string(output))
	if containerID == "" {
		err := fmt.Errorf("empty container ID received")
		logger.Err("Creating container", err)
		ws.EventEmitter.Error("Creating container", err)
		return nil, err
	}
	logger.ContainerStart(containerID)
	ws.EventEmitter.ContainerStart(containerID)

	node := node.Node{
		ContainerID: containerID,
		Address:     fmt.Sprintf("http://localhost:%d", externalPort),
		Metrics: node.NodeMetrics{
			Health: "unknown",
		},
	}

	logger.Info(fmt.Sprintf("Started server @ http://localhost:%d", externalPort))
	ws.EventEmitter.Info(fmt.Sprintf("Started server @ http://localhost:%d", externalPort))
	return &node, nil
}

func StartContainer(externalPort int, dockerInfo *config.DockerConfig) (*node.Node, error) {
	cli, err := createDockerClient()
	if err != nil {
		logger.Err("Failed to create Docker client", err)
		return nil, err
	}

	ctx := context.Background()
	internalPort := fmt.Sprintf("%d/tcp", dockerInfo.InternalPort)

	containerConfig := &container.Config{
		Image: dockerInfo.Image,
		ExposedPorts: nat.PortSet{
			nat.Port(internalPort): {},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(internalPort): {
				{HostPort: fmt.Sprintf("%d", externalPort)},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		logger.Err("Failed to create container", err)
		ws.EventEmitter.Error("Failed to create container", err)
		return nil, err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		logger.Err("Failed to start container", err)
		ws.EventEmitter.Error("Failed to start container", err)
		return nil, err
	}

	logger.ContainerStart(resp.ID)
	ws.EventEmitter.ContainerStart(resp.ID)

	node := node.Node{
		ContainerID: resp.ID,
		Address:     fmt.Sprintf("http://localhost:%d", externalPort),
		Metrics: node.NodeMetrics{
			Health: "unknown",
		},
	}

	logger.Info(fmt.Sprintf("Started server @ http://localhost:%d", externalPort))
	ws.EventEmitter.Info(fmt.Sprintf("Started server @ http://localhost:%d", externalPort))
	return &node, nil
}

func createDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost("unix:///var/run/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
