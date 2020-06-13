package central

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

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

// DockerClient is a wrapper over the docker cllient struct.
type DockerClient struct {
	c *client.Client
}

// createAndStartContainer creates and starts a new container
func (cli *DockerClient) createAndStartContainer(ctx context.Context, containerConfig *Container) (string, error) {
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
	}

	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPortRule: []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", containerConfig.ContainerPort)}},
		},
	}

	createdConfig, err := cli.c.ContainerCreate(ctx, containerConf, hostConf, nil, containerConfig.Name)
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

// stopAndRemoveContainer stops and removes the container from the docker host
func (cli *DockerClient) stopAndRemoveContainer(ctx context.Context, containerID string) error {
	err := cli.c.ContainerStop(ctx, containerID, nil)
	if err != nil {
		return err
	}

	return cli.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

// inspectContainer returns container state by using containerID
func (cli *DockerClient) inspectContainer(ctx context.Context, containerID string) (*types.ContainerState, error) {
	containerStatus, err := cli.c.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return containerStatus.ContainerJSONBase.State, nil
}

// listContainers returns a list of containers ids
func (cli *DockerClient) listContainers(ctx context.Context) ([]types.Container, error) {
	containers, err := cli.c.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// readLogs returns the logs from a docker container
func (cli *DockerClient) readLogs(ctx context.Context, containerID, tail string) ([]string, error) {
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

// pullImage requests the host to pull an image
func (cli *DockerClient) pullImage(ctx context.Context, image string) error {
	out, err := cli.c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close() //nolint:errcheck

	out2 := strings.NewReader("")
	_, err = io.Copy(os.Stdout, out2)
	if err != nil {
		return err
	}
	return nil
}
