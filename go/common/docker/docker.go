package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func StartNewContainer(containerName, image string, cmds []string, ports []int, envs map[string]string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// Check if the image exists locally
	_, _, err = cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		// unexpected error
		if !client.IsErrNotFound(err) {
			return "", err
		}

		// Pull the image from remote
		fmt.Printf("Image %s not found locally. Pulling from remote...\n", image)
		pullReader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		defer pullReader.Close()
		go func() {
			_, err = io.Copy(os.Stdout, pullReader)
			if err != nil {
				fmt.Println(err)
			}
		}()

		// Wait until the image is available in the local Docker image cache
		imageFilter := filters.NewArgs()
		imageFilter.Add("reference", image)
		imageListOptions := types.ImageListOptions{Filters: imageFilter}
		for {
			imageSummaries, err := cli.ImageList(context.Background(), imageListOptions)
			if err != nil {
				return "", err
			}
			if len(imageSummaries) > 0 {
				break
			}
		}

		// Image is available
		fmt.Printf("Image %s is available!\n", image)
	} else {
		fmt.Printf("Image %s found locally.\n", image)
	}

	// convert env vars
	envVars := make([]string, 0, len(envs))
	for k, v := range envs {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	// expose ports
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, port := range ports {
		portBindings[nat.Port(fmt.Sprintf("%d/tcp", port))] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", port)}}
		exposedPorts[nat.Port(fmt.Sprintf("%d/tcp", port))] = struct{}{}
	}

	// create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Entrypoint:   cmds,
		Tty:          false,
		ExposedPorts: exposedPorts,
		Env:          envVars,
	},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"node_network": {
					NetworkID: "node_network",
				},
			},
		}, containerName)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		return "", err
	}

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return resp.ID, nil
}

func WaitForContainerToFinish(containerID string, timeout time.Duration) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	// Wait for the container to finish with a timeout of one minute
	statusCh, errCh := cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with non-zero status code")
		}
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for container to finish")
	}

	return nil
}
