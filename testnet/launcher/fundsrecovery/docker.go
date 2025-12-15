package fundsrecovery

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
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
	fmt.Printf("Starting funds recovery with config: %s\n", n.cfg)
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"--network",
		"layer1",
		"scripts/recover/recover_testnet_funds.ts",
		"--verbose",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`
		{
			"layer1": {
				"url": "%s",
				"useGateway": false,
				"live": false,
				"saveDeployments": true,
				"accounts": ["%s"]
			}
		}`, n.cfg.l1HTTPURL, n.cfg.l1privateKey),
		"NETWORK_CONFIG_ADDR": n.cfg.networkConfigAddress,
		"RECEIVER_ADDRESS":    n.cfg.receiverAddress,
	}

	// Mount only the scripts directory to use updated scripts
	// This avoids conflicts with node_modules which are platform-specific
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	scriptsPath := filepath.Join(cwd, "contracts", "scripts")
	fmt.Printf("Mounting local scripts directory: %s\n", scriptsPath)
	volumes := map[string]string{
		scriptsPath: "/home/obscuro/go-obscuro/contracts/scripts",
	}

	containerID, err := docker.StartNewContainer("recover-funds", n.cfg.dockerImage, cmds, nil, envs, nil, volumes, false)
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
