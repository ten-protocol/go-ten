package l1upgrade

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
	"github.com/ten-protocol/go-ten/go/common/docker"
)

type UpgradeContracts struct {
	cfg         *Config
	containerID string
}

func NewUpgradeContracts(cfg *Config) (*UpgradeContracts, error) {
	return &UpgradeContracts{
		cfg: cfg,
	}, nil
}

func (s *UpgradeContracts) Start() error {
	fmt.Printf("Starting upgrade conracts with config: %s\n", s.cfg)
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"--network",
		"layer1",
		"scripts/upgrade/001_upgrade_contracts.ts",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`{ 
            "layer1": {
                "url": "%s",
                "live": false,
                "saveDeployments": true,
                "accounts": [ "%s" ]
            }
        }`, s.cfg.l1HTTPURL, s.cfg.privateKey),
		"NETWORK_CONFIG_ADDR": s.cfg.networkConfigAddress,
	}

	fmt.Printf("Starting upgrade contracts script. NeworkConfigAddress: %s", s.cfg.networkConfigAddress)

	containerID, err := docker.StartNewContainer(
		"upgrade-contracts",
		s.cfg.dockerImage,
		cmds,
		nil,
		envs,
		nil,
		nil,
		false,
	)
	if err != nil {
		return err
	}
	s.containerID = containerID
	return nil
}

func (s *UpgradeContracts) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	// make sure the container has finished execution
	err = docker.WaitForContainerToFinish(s.containerID, 15*time.Minute)
	if err != nil {
		fmt.Println("Error waiting for container to finish: ", err)
		s.PrintLogs(cli)
		return err
	}

	return nil
}

func (s *UpgradeContracts) PrintLogs(cli *client.Client) {
	logsOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), s.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", s.containerID, err)
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", s.containerID)
	}
	fmt.Println(buf.String())
}
