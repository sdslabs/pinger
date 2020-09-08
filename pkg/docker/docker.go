// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package docker

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

// ContainerConf is the configuration for creating a container.
type ContainerConf struct {
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

// Client is a wrapper over the docker client struct.
type Client struct {
	c *client.Client
}

// NewClient creates a new docker client.
func NewClient() (Client, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return Client{}, err
	}

	return Client{cli}, nil
}

// CreateAndStartContainer creates and starts a new container.
func (cli *Client) CreateAndStartContainer(ctx context.Context, containerConfig *ContainerConf) (string, error) {
	// convert map to list of strings
	envArr := make([]string, len(containerConfig.Env))
	i := 0
	for key, value := range containerConfig.Env {
		envArr[i] = fmt.Sprintf("%s=%v", key, value)
		i++
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

// StopAndRemoveContainer stops and removes the container from the docker host.
func (cli *Client) StopAndRemoveContainer(ctx context.Context, containerID string) error {
	err := cli.c.ContainerStop(ctx, containerID, nil)
	if err != nil {
		return err
	}

	return cli.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

// InspectContainer returns container state by using containerID.
func (cli *Client) InspectContainer(ctx context.Context, containerID string) (*types.ContainerState, error) {
	containerStatus, err := cli.c.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	return containerStatus.ContainerJSONBase.State, nil
}

// ListContainers returns a list of containers.
func (cli *Client) ListContainers(ctx context.Context) ([]types.Container, error) {
	containers, err := cli.c.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// ReadLogs returns the logs from a docker container.
func (cli *Client) ReadLogs(ctx context.Context, containerID, tail string) ([]string, error) {
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

// PullImage requests the host to pull an image.
func (cli *Client) PullImage(ctx context.Context, image string) error {
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
