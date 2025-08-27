package transferunpauserroles

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

type RoleTransfer struct {
	cfg         *Config
	containerID string
}

func NewRoleTransfer(cfg *Config) (*RoleTransfer, error) {
	return &RoleTransfer{
		cfg: cfg,
	}, nil
}

func (r *RoleTransfer) Start() error {
	fmt.Printf("Starting unpauser role transfer with config: %s\n", r.cfg)
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"scripts/pause/001_transfer_unpauser_role.ts",
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
		}`, r.cfg.l1HTTPURL, r.cfg.privateKey),
		"NETWORK_CONFIG_ADDR":          r.cfg.networkConfigAddr,
		"MULTISIG_ADDR":                r.cfg.multisigAddr,
		"MERKLE_TREE_MESSAGE_BUS_ADDR": r.cfg.merkleMessageBusAddr,
	}

	fmt.Printf("Starting unpauser role transfer script. NetworkConfigAddress: %s, MultisigAddress: %s\n",
		r.cfg.networkConfigAddr, r.cfg.multisigAddr)

	containerID, err := docker.StartNewContainer(
		"transfer-unpauser-roles",
		r.cfg.dockerImage,
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
	r.containerID = containerID

	return nil
}

func (r *RoleTransfer) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = docker.WaitForContainerToFinish(r.containerID, 15*time.Minute)
	if err != nil {
		r.PrintLogs(cli)
		return err
	}

	return nil
}

func (r *RoleTransfer) PrintLogs(cli *client.Client) {
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
	out, err := cli.ContainerLogs(ctx, r.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", r.containerID, err)
		return
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", r.containerID)
		return
	}
	fmt.Println(buf.String())
}
