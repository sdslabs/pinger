package central

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

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

	err = cli.c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return containerID, nil
}

// StopAndRemoveContainer stops and removes the container from the docker host
func (cli DockerClient) StopAndRemoveContainer(ctx context.Context, containerID string) error {
	inspectVal, _, err := cli.c.ContainerInspectWithRaw(ctx, containerID, false)
	if err != nil {
		return err
	}

	if inspectVal.ID != "" {
		err = cli.c.ContainerStop(ctx, containerID, nil)
		if err != nil {
			return err
		}

		err = cli.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}

	return nil
}

// InspectContainer returns container state by using containerID
func (cli DockerClient) InspectContainer(ctx context.Context, containerID string) (*types.ContainerState, error) {
	containerStatus, err := cli.c.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return containerStatus.ContainerJSONBase.State, nil
}

// ListContainers returns a list of containers ids
func (cli DockerClient) ListContainers(ctx context.Context) ([]types.Container, error) {
	containers, err := cli.c.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	list := make([]types.Container, 1)

	for _, container := range containers {
		list = append(list, container)
	}
	return list, nil
}

// ReadLogs returns the logs from a docker container
func (cli DockerClient) ReadLogs(ctx context.Context, containerID, tail string) ([]string, error) {
	reader, err := cli.c.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: true,
		Tail:       tail,
	})

	if err != nil {
		return nil, err
	}

	defer reader.Close()

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
			return logs, err
		}
	}
}

// PullImage requests the host to pull an image
func (cli DockerClient) PullImage(ctx context.Context, image string) error {
	out, err := cli.c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)
	return nil
}

// BuildImage sends request to build image
func (cli DockerClient) BuildImage(ctx context.Context, tarCtxPath, dockerCtxFile, tagName string) error {
	buildContext, err := os.Open(tarCtxPath)
	if err != nil {
		fmt.Errorf("Error while opening staged file :: %s", tarCtxPath)
	}
	defer buildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{tagName},
		NoCache:    true,
		Remove:     true,
		Dockerfile: dockerCtxFile,
	}

	image, err := cli.c.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		log.Println(err, ":unable to build image")
	}
	return nil
}

// RemoveImage removes an image by imageID
func (cli DockerClient) RemoveImage(ctx context.Context, imageID string) error {
	inspectVal, _, err := cli.c.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return err
	}

	if inspectVal.ID != "" {
		_, err = cli.c.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{})
	}

	return nil
}
