package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	dockerimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

const _networkName = "node_network"

// StartNewContainer - volumes is a map from volume name to dir name it will have within the container. If a volume doesn't exist this will create it.
func StartNewContainer(containerName, image string, cmds []string, ports []int, envs, devices, volumes map[string]string, autoRestart bool) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// Check if the image exists locally
	_, err = cli.ImageInspect(context.Background(), image)
	if err != nil {
		// unexpected error
		if !cerrdefs.IsNotFound(err) {
			return "", err
		}

		err = waitAndPullRemoteImage(image, cli)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Printf("Image %s found locally.\n", image)
	}

	// Check if the network already exists
	err = createNetwork(_networkName, cli)
	if err != nil {
		return "", err
	}

	// convert devices
	deviceMapping := make([]container.DeviceMapping, 0, len(devices))
	for k, v := range devices {
		deviceMapping = append(deviceMapping, container.DeviceMapping{
			PathOnHost:        k,
			PathInContainer:   v,
			CgroupPermissions: "rwm",
		})
	}

	mountVolumes := make([]mount.Mount, 0, len(volumes))
	for v, mntTarget := range volumes {
		vol, err := ensureVolumeExists(cli, v)
		if err != nil {
			return "", err
		}
		mountVolumes = append(mountVolumes, mount.Mount{
			Type:   mount.TypeVolume,
			Source: vol.Name,
			Target: mntTarget,
		})
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

	// set log rotations
	logOptions := map[string]string{
		"max-size": "10m",
		"max-file": "3",
	}

	hc := container.HostConfig{
		PortBindings: portBindings,
		Mounts:       mountVolumes,
		Resources:    container.Resources{Devices: deviceMapping},
		LogConfig:    container.LogConfig{Type: "json-file", Config: logOptions},
	}

	if autoRestart {
		hc.RestartPolicy = container.RestartPolicy{Name: "unless-stopped"}
	}

	// create the container
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        image,
			Entrypoint:   cmds,
			Tty:          false,
			ExposedPorts: exposedPorts,
			Env:          envVars,
		},
		&hc,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				_networkName: {
					NetworkID: _networkName,
				},
			},
		},
		nil,
		containerName,
	)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		return "", err
	}

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return resp.ID, nil
}

func StopAndRemove(containerName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, containerName, container.StopOptions{})
	if err != nil {
		return err
	}

	return cli.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
}

func ensureVolumeExists(cli *client.Client, volumeName string) (*volume.Volume, error) {
	ctx := context.Background()
	allVolumes, err := cli.VolumeList(ctx, volume.ListOptions{Filters: filters.NewArgs()})
	if err != nil {
		return nil, fmt.Errorf("unable to list volumes - %w", err)
	}
	for _, v := range allVolumes.Volumes {
		if v.Name == volumeName {
			fmt.Printf("Volume %s found - reusing existing volume! \n", volumeName)
			return v, nil
		}
	}
	// volume doesn't exist, so create it
	vol, err := cli.VolumeCreate(ctx, volume.CreateOptions{
		Driver: "local",
		Name:   volumeName,
	})
	fmt.Println("Volume not found in docker, created: ", volumeName)
	return &vol, err
}

func createNetwork(networkName string, cli *client.Client) error {
	// Check if the network already exists
	networkFilter := network.ListOptions{Filters: filters.NewArgs()}
	networkFilter.Filters.Add("name", networkName)
	existingNetworks, err := cli.NetworkList(context.Background(), networkFilter)
	if err != nil {
		return err
	}

	if len(existingNetworks) == 0 {
		// Create the network if it doesn't exist
		_, err = cli.NetworkCreate(
			context.Background(),
			networkName,
			network.CreateOptions{
				Driver:     "bridge",
				Attachable: true,
				Ingress:    false,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func waitAndPullRemoteImage(image string, cli *client.Client) error {
	// Pull the image from remote
	fmt.Printf("Image %s not found locally. Pulling from remote...\n", image)
	pullReader, err := cli.ImagePull(context.Background(), image, dockerimage.PullOptions{})
	if err != nil {
		return err
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
	imageListOptions := dockerimage.ListOptions{Filters: imageFilter}
	for {
		imageSummaries, err := cli.ImageList(context.Background(), imageListOptions)
		if err != nil {
			return err
		}
		if len(imageSummaries) > 0 {
			break
		}
	}

	// Image is available
	fmt.Printf("Image %s is available!\n", image)
	return nil
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
	case err = <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with non-zero status code - Status Code: %d", status.StatusCode)
		}
	case <-time.After(timeout):
		return fmt.Errorf("timeout after %s waiting for container to finish", timeout)
	}

	return nil
}

// ExecInContainer executes a command in a running container and returns an error if it fails
func ExecInContainer(containerName string, cmd []string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: false,
		AttachStderr: false,
	}

	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return err
	}

	err = cli.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}

	// wait for the exec to finish
	inspectResp, err := cli.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return err
	}

	for inspectResp.Running {
		time.Sleep(100 * time.Millisecond)
		inspectResp, err = cli.ContainerExecInspect(ctx, execID.ID)
		if err != nil {
			return err
		}
	}

	if inspectResp.ExitCode != 0 {
		return fmt.Errorf("command exited with code %d", inspectResp.ExitCode)
	}

	return nil
}
