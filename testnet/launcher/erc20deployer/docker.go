package erc20deployer

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

type ERC20Deployer struct {
	cfg         *Config
	containerID string
}

func NewERC20Deployer(cfg *Config) (*ERC20Deployer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &ERC20Deployer{
		cfg: cfg,
	}, nil
}

func (d *ERC20Deployer) Start() error {
	fmt.Printf("Starting ERC20 token deployment with config: %s\n", d.cfg)

	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"scripts/deploy/deploy_erc20_token.ts",
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
		}`, d.cfg.l1HTTPURL, d.cfg.privateKey),
		"TOKEN_NAME":          d.cfg.tokenName,
		"TOKEN_SYMBOL":        d.cfg.tokenSymbol,
		"TOKEN_DECIMALS":      d.cfg.tokenDecimals,
		"TOKEN_SUPPLY":        d.cfg.tokenSupply,
		"NETWORK_CONFIG_ADDR": d.cfg.networkConfigAddr,
	}

	// Mount scripts and src directories so contracts can be compiled
	// This avoids conflicts with node_modules which are platform-specific
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	scriptsPath := filepath.Join(cwd, "contracts", "scripts")
	srcPath := filepath.Join(cwd, "contracts", "src")
	volumes := map[string]string{
		scriptsPath: "/home/obscuro/go-obscuro/contracts/scripts",
		srcPath:     "/home/obscuro/go-obscuro/contracts/src",
	}

	fmt.Println("About to call docker.StartNewContainer...")
	containerID, err := docker.StartNewContainer(
		"deploy-erc20-token",
		d.cfg.dockerImage,
		cmds,
		nil,
		envs,
		nil,
		volumes,
		false,
	)
	if err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		return err
	}
	fmt.Printf("Container started with ID: %s\n", containerID)
	d.containerID = containerID

	return nil
}

func (d *ERC20Deployer) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = docker.WaitForContainerToFinish(d.containerID, 10*time.Minute)
	if err != nil {
		d.PrintLogs(cli)
		return err
	}

	// Print logs on success too
	d.PrintLogs(cli)

	return nil
}

func (d *ERC20Deployer) PrintLogs(cli *client.Client) {
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
	out, err := cli.ContainerLogs(ctx, d.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", d.containerID, err)
		return
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", d.containerID)
		return
	}
	fmt.Println(buf.String())
}
