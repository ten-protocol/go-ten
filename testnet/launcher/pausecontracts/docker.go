package pausecontracts

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

type PauseAllContracts struct {
	cfg         *Config
	containerID string
}

func NewPauseAllContracts(cfg *Config) (*PauseAllContracts, error) {
	return &PauseAllContracts{
		cfg: cfg,
	}, nil
}

func (p *PauseAllContracts) Start() error {
	fmt.Printf("Starting pause all contracts with config: %s\n", p.cfg)
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"scripts/pause/002_pause_all_contracts.ts",
		"--network",
		"layer1",
		"--verbose",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`
		{
			"layer1": {
				"url": "%s",
				"live": false,
				"saveDeployments": true,
				"accounts": ["%s"]
			}
		}`, p.cfg.l1HTTPURL, p.cfg.privateKey),
		"NETWORK_CONFIG_ADDR":          p.cfg.networkConfigAddr,
		"MERKLE_TREE_MESSAGE_BUS_ADDR": p.cfg.merkleMessageBusAddr,
		"ACTION":                       p.cfg.action,
	}

	fmt.Printf("Starting pause all contracts script. NetworkConfigAddress: %s\n",
		p.cfg.networkConfigAddr)

	containerID, err := docker.StartNewContainer(
		"pause-all-contracts",
		p.cfg.dockerImage,
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
	p.containerID = containerID

	return nil
}

func (p *PauseAllContracts) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = docker.WaitForContainerToFinish(p.containerID, 15*time.Minute)
	if err != nil {
		p.PrintLogs(cli)
		return err
	}

	return nil
}

func (p *PauseAllContracts) PrintLogs(cli *client.Client) {
	if cli == nil {
		fmt.Println("Docker client is nil, cannot print logs")
		return
	}

	logsOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Read the container logs
	out, err := cli.ContainerLogs(ctx, p.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", p.containerID, err)
		return
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", p.containerID)
		return
	}
	fmt.Println(buf.String())
}
