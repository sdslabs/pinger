package central

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Container is the configuration for creating a container
type Container struct {
	// Name of the container
	Name string
	// Image used for creating the container
	Image string
	// Port on which container is running
	ContainerPort int
	// Port of the docker container in the host system
	HostPort int
	// Environment variables
	Env map[string]interface{}
}

type DockerClient struct {
	c *client.Client
}

// CreateAndStartContainer creates and starts a new container
func (cli DockerClient) CreateAndStartContainer(ctx context.Context, containerConfig *Container) (string, error) {
	// convert map to list of strings
	envArr := []string{}
	for key, value := range containerConfig.Env {
		envArr = append(envArr, fmt.Sprintf("%s=%v", key, value))
	}

	containerPortRule := nat.Port(fmt.Sprintf(`%d/tcp`, containerConfig.ContainerPort))

	containerConf := &container.Config{
		Image: containerConfig.Image,
		Env:   envArr,
		ExposedPorts: nat.PortSet{
			containerPortRule: struct{}{},
		},
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD-SHELL", fmt.Sprintf("curl --fail --silent http://localhost:%d/ || exit 1", containerConfig.ContainerPort)},
			Interval: 5 * time.Second,
			Timeout:  10 * time.Second,
			Retries:  3,
		},
	}

	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(containerPortRule): []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", containerConfig.ContainerPort)}},
		},
	}

	createdConfig, err := cli.c.ContainerCreate(ctx, containerConf, hostConf, nil, containerConfig.Name)
	if err != nil {
		return "", err
	}
	containerID := createdConfig.ID
	if err := cli.c.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return containerID, nil
}

// StopAndRemoveContainer stops and removes the container from the docker host
func (cli DockerClient) StopAndRemoveContainer(ctx context.Context, containerID string) error {
	return cli.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveLinks: false,
		Force:       true,
	})
}
