package node

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/network"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerNode struct {
	cfg *Config
}

func NewDockerNode(cfg *Config) (*DockerNode, error) {
	return &DockerNode{
		cfg: cfg,
	}, nil // todo: add config validation
}

func (d *DockerNode) Start() error {
	err := d.startEnclave()
	if err != nil {
		return err
	}

	err = d.startHost()
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerNode) startHost() error {
	cmd := []string{
		"/home/obscuro/go-obscuro/go/host/main/main",
		"-l1NodeHost", d.cfg.l1Host,
		"-l1NodePort", fmt.Sprintf("%d", d.cfg.l1WSPort),
		"-enclaveRPCAddress", fmt.Sprintf("enclave:%d", d.cfg.enclaveHTTPPort),
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-privateKey", d.cfg.privateKey,
		"-clientRPCHost", "0.0.0.0",
		"-logPath", "sys_out",
		"-logLevel", "4",
		"-isGenesis", fmt.Sprintf("%t", d.cfg.isGenesis),
		"-nodeType", fmt.Sprintf("%s", d.cfg.nodeType),
		"-profilerEnabled", "false",
		"-p2pPublicAddress", fmt.Sprintf("0.0.0.0:%d", d.cfg.hostP2PPort),
	}

	return startNewContainer("host", d.cfg.hostImage, cmd)

}

func (d *DockerNode) startEnclave() error {
	cmd := []string{
		"ego", "run", "/home/obscuro/go-obscuro/go/enclave/main/main",
		"-hostID", d.cfg.hostID,
		"-address", fmt.Sprintf("0.0.0.0:%d", d.cfg.enclaveHTTPPort),
		"-nodeType", d.cfg.nodeType,
		"-useInMemoryDB", "false",
		"-sqliteDBPath", "/data/sqlite.db",
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-hostAddress", fmt.Sprintf("host:%d", d.cfg.hostHTTPPort),
		"-profilerEnabled", "false",
		"-hostAddress", fmt.Sprintf("host:%d", d.cfg.hostP2PPort),
		"-logPath", "sys_out",
		"-logLevel", "2",
		"-sequencerID", d.cfg.sequencerID,
		"-messageBusAddress", d.cfg.messageBusContractAddress,
	}

	return startNewContainer("enclave", d.cfg.enclaveImage, cmd)
}

func startNewContainer(containerName, image string, cmds []string) error {
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

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      image,
		Entrypoint: cmds,
		Tty:        false,
	}, nil, &network.NetworkingConfig{
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
