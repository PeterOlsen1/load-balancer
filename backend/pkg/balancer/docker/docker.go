package docker

import (
	"context"
	"fmt"
	"load-balancer/pkg/balancer/node"
	"load-balancer/pkg/config"
	"load-balancer/pkg/logger"
	"load-balancer/pkg/ws"

	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

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

func StopContainer(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is empty")
	}

	cli, err := createDockerClient()
	if err != nil {
		logger.Err("Failed to create Docker client", err)
		ws.EventEmitter.Error("Failed to create Docker client", err)
		return err
	}

	ctx := context.Background()
	timeout := 1

	err = cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to stop container %s", containerID), err)
		ws.EventEmitter.Error(fmt.Sprintf("Failed to stop container %s", containerID), err)
		return err
	}

	err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{})
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to remove container %s", containerID), err)
		ws.EventEmitter.Error(fmt.Sprintf("Failed to remove container %s", containerID), err)
		return err
	}

	logger.ContainerStop(containerID)
	ws.EventEmitter.ContainerStop(containerID)

	return nil
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
