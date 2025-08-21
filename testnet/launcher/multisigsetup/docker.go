package multisigsetup

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

type MultisigSetup struct {
	cfg         *Config
	containerID string
}

func NewMultisigSetup(cfg *Config) (*MultisigSetup, error) {
	return &MultisigSetup{
		cfg: cfg,
	}, nil
}

func (s *MultisigSetup) Start() error {
	fmt.Printf("Starting multisig setup with config: %s\n", s.cfg)
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"scripts/multisig/phase_1/001_direct_multisig_setup.ts",
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
		}`, s.cfg.l1HTTPURL, s.cfg.privateKey),
		"NETWORK_CONFIG_ADDR": s.cfg.networkConfigAddress,
		"MULTISIG_ADDR":       s.cfg.multisigAddress,
		"PROXY_ADMIN_ADDR":    s.cfg.proxyAdminAddress,
	}

	fmt.Printf("Starting multisig setup script. NetworkConfigAddress: %s, MultisigAddress: %s\n",
		s.cfg.networkConfigAddress, s.cfg.multisigAddress)

	containerID, err := docker.StartNewContainer(
		"multisig-setup",
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

func (s *MultisigSetup) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = docker.WaitForContainerToFinish(s.containerID, 15*time.Minute)
	if err != nil {
		s.PrintLogs(cli)
		return err
	}

	return nil
}

func (s *MultisigSetup) PrintLogs(cli *client.Client) {
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
	out, err := cli.ContainerLogs(ctx, s.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", s.containerID, err)
		return
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", s.containerID)
		return
	}
	fmt.Println(buf.String())
}
