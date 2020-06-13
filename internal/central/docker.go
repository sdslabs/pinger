package central

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// containerConf is the configuration for creating a container.
type containerConf struct {
	// Name of the container
	name string
	// Image used for creating the container
	image string
	// Port on which container is running
	containerPort int
	// Port of the docker container in the host system
	hostPort int
	// Environment variables
	env map[string]interface{}
}

// dockerClient is a wrapper over the docker cllient struct.
type dockerClient struct {
	c *client.Client
}

// newDockerClient creates a new docker client.
func newDockerClient() (*dockerClient, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &dockerClient{cli}, nil
}

// createAndStartContainer creates and starts a new container.
func (cli *dockerClient) createAndStartContainer(ctx context.Context, containerConfig *containerConf) (string, error) {
	// convert map to list of strings
	envArr := make([]string, len(containerConfig.env))
	i := 0
	for key, value := range containerConfig.env {
		envArr[i] = fmt.Sprintf("%s=%v", key, value)
		i++
	}

	containerPortRule := nat.Port(fmt.Sprintf(`%d/tcp`, containerConfig.containerPort))

	containerConf := &container.Config{
		Image: containerConfig.image,
		Env:   envArr,
		ExposedPorts: nat.PortSet{
			containerPortRule: struct{}{},
		},
	}

	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPortRule: []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", containerConfig.containerPort)}},
		},
	}

	createdConfig, err := cli.c.ContainerCreate(ctx, containerConf, hostConf, nil, containerConfig.name)
	if err != nil {
		return "", err
	}
	containerID := createdConfig.ID

	err = cli.c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return containerID, nil
}

// stopAndRemoveContainer stops and removes the container from the docker host.
func (cli *dockerClient) stopAndRemoveContainer(ctx context.Context, containerID string) error {
	err := cli.c.ContainerStop(ctx, containerID, nil)
	if err != nil {
		return err
	}

	return cli.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

// inspectContainer returns container state by using containerID.
func (cli *dockerClient) inspectContainer(ctx context.Context, containerID string) (*types.ContainerState, error) {
	containerStatus, err := cli.c.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return containerStatus.ContainerJSONBase.State, nil
}

// listContainers returns a list of containers.
func (cli *dockerClient) listContainers(ctx context.Context) ([]types.Container, error) {
	containers, err := cli.c.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// readLogs returns the logs from a docker container.
func (cli *dockerClient) readLogs(ctx context.Context, containerID, tail string) ([]string, error) {
	reader, err := cli.c.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: true,
		Tail:       tail,
	})

	if err != nil {
		return nil, err
	}

	defer reader.Close() //nolint:errcheck

	logs := []string{}
	hdr := make([]byte, 8)

	for {
		_, err = reader.Read(hdr)

		if err != nil {
			return logs, err
		}

		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = reader.Read(dat)
		logs = append(logs, string(dat))

		if err != nil {
			if err == io.EOF {
				break
			}
			return logs, err
		}
	}
	return logs, nil
}

// pullImage requests the host to pull an image.
func (cli *dockerClient) pullImage(ctx context.Context, image string) error {
	out, err := cli.c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close() //nolint:errcheck

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, out)
	if err != nil {
		return err
	}
	return nil
}
