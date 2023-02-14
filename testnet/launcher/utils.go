package launcher

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func startNewContainer(containerName, image string, cmds []string, ports []int) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	// Check if the image exists locally
	_, _, err = cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		if client.IsErrNotFound(err) {
			fmt.Printf("Image %s not found locally. Pulling from remote...\n", image)
			// Pull the image from remote
			pullReader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}
			defer pullReader.Close()
			go func() {
				_, err = io.Copy(os.Stdout, pullReader)
				if err != nil {
					fmt.Println(err)
				}
			}()
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("Image %s found locally.\n", image)
	}

	exposedPort := nat.PortMap{}
	for _, port := range ports {
		exposedPort[nat.Port(fmt.Sprintf("%d", port))] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", port)}}
	}

	hostConfig := &container.HostConfig{
		PortBindings: exposedPort,
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      image,
		Entrypoint: cmds,
		Tty:        false,
	},
		hostConfig,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"node_network": {
					NetworkID: "node_network",
				},
			},
		}, containerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	//
	//statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	//select {
	//case err := <-errCh:
	//	if err != nil {
	//		panic(err)
	//	}
	//case <-statusCh:
	//}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
