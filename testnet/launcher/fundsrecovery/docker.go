package fundsrecovery

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/docker"
)

type FundsRecovery struct {
	cfg         *Config
	containerID string
}

func NewFundsRecovery(cfg *Config) (*FundsRecovery, error) {
	return &FundsRecovery{
		cfg: cfg,
	}, nil
}

func (n *FundsRecovery) Start() error {
	fmt.Printf("Starting L2 contract deployer with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{
		"npx", "hardhat", "deploy",
		"--network", "layer1",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`
{
        "layer1" : {
            "url" : "%s",
            "useGateway" : false,
            "live" : false,
            "saveDeployments" : true,
            "deploy": [
                "deployment_scripts/testnet/recoverfunds"
            ],
            "accounts": [
                "%s"
            ]
        }
    }
`, n.cfg.l1HTTPURL, n.cfg.l1privateKey),
	}

	containerID, err := docker.StartNewContainer("recover-funds", n.cfg.dockerImage, cmds, nil, envs, nil, nil, false)
	if err != nil {
		return err
	}
	n.containerID = containerID
	return nil
}

func (n *FundsRecovery) GetID() string {
	return n.containerID
}

func (n *FundsRecovery) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	// make sure the container has finished execution
	err = docker.WaitForContainerToFinish(n.containerID, 10*time.Minute)
	if err != nil {
		n.PrintLogs(cli)
		return err
	}

	// if we want to read anything from the container logs we can do it here (see RetrieveL1ContractAddresses as example)

	return nil
}

func (n *FundsRecovery) PrintLogs(cli *client.Client) {
	logsOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), n.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", n.containerID, err)
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", n.containerID)
	}
	fmt.Println(buf.String())
}
