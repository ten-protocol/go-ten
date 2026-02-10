package bridgetokenwhitelist

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

type BridgeTokenWhitelister struct {
	cfg         *Config
	containerID string
}

func NewBridgeTokenWhitelister(cfg *Config) (*BridgeTokenWhitelister, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &BridgeTokenWhitelister{
		cfg: cfg,
	}, nil
}

func (w *BridgeTokenWhitelister) Start() error {
	fmt.Printf("Starting bridge token whitelisting with config: %s\n", w.cfg)

	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"scripts/bridge/whitelist_and_register_tokens.ts",
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
		}`, w.cfg.l1HTTPURL, w.cfg.privateKey),
		"TOKEN_ADDRESS":       w.cfg.tokenAddress,
		"TOKEN_NAME":          w.cfg.tokenName,
		"TOKEN_SYMBOL":        w.cfg.tokenSymbol,
		"NETWORK_CONFIG_ADDR": w.cfg.networkConfigAddr,
		"L2_RPC_URL":          w.cfg.l2RPCURL,
	}

	// Mount scripts and src directories so contracts can be compiled if needed
	// This avoids conflicts with node_modules which are platform-specific
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	scriptsPath := filepath.Join(cwd, "contracts", "scripts")
	volumes := map[string]string{
		scriptsPath: "/home/obscuro/go-obscuro/contracts/scripts",
	}

	containerID, err := docker.StartNewContainer(
		"whitelist-bridge-token",
		w.cfg.dockerImage,
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
	w.containerID = containerID

	return nil
}

func (w *BridgeTokenWhitelister) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	err = docker.WaitForContainerToFinish(w.containerID, 10*time.Minute)
	if err != nil {
		w.PrintLogs(cli)
		return err
	}

	// Print logs on success too
	w.PrintLogs(cli)

	return nil
}

func (w *BridgeTokenWhitelister) PrintLogs(cli *client.Client) {
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
	out, err := cli.ContainerLogs(ctx, w.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", w.containerID, err)
		return
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", w.containerID)
		return
	}
	fmt.Println(buf.String())
}
