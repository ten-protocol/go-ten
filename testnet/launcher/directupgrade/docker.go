package directupgrade

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

type DirectUpgrade struct {
	cfg         *Config
	containerID string
}

func NewDirectUpgrade(cfg *Config) (*DirectUpgrade, error) {
	return &DirectUpgrade{
		cfg: cfg,
	}, nil
}

func (s *DirectUpgrade) Start() error {
	fmt.Printf("Starting direct upgrade with config: %s\n", s.cfg)
	var err error
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"scripts/multisig/phase_1/002_direct_upgrade.ts",
		"--network",
		"layer1",
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
		}`, s.cfg.l1HTTPURL, s.cfg.privateKey),
		"NETWORK_CONFIG_ADDR": s.cfg.networkConfigAddress,
		"MULTISIG_ADDR":       s.cfg.multisigAddress,
		"PROXY_ADMIN_ADDR":    s.cfg.proxyAdminAddress,
	}

	fmt.Printf("Starting direct upgrade script. NetworkConfigAddress: %s, MultisigAddress: %s, ProxyAdminAddress: %s\n",
		s.cfg.networkConfigAddress, s.cfg.multisigAddress, s.cfg.proxyAdminAddress)

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

	containerID, err := docker.StartNewContainer(
		"direct-upgrade",
		s.cfg.dockerImage,
		cmds,
		nil,
		envs,
		nil,
		volumes,
		false,
	)
	if err != nil {
		return err
	}
	s.containerID = containerID

	return nil
}

func (s *DirectUpgrade) WaitForFinish() error {
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

func (s *DirectUpgrade) PrintLogs(cli *client.Client) {
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
