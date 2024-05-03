package l1contractdeployer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/docker"
)

type ContractDeployer struct {
	cfg         *config.L1ContractDeployerConfig
	netCfg      *config.NetworkInputConfig
	containerID string
}

func NewDockerContractDeployer(cfg *config.L1ContractDeployerConfig) (*ContractDeployer, error) {
	return &ContractDeployer{
		cfg: cfg,
	}, nil // todo (@pedro) - add validation
}

func NewDockerContractDeployerFromTestnetConfig(cfg *config.TestnetConfig) (*ContractDeployer, error) {
	return &ContractDeployer{
		cfg:    &cfg.L1ContractDeployer,
		netCfg: &cfg.Network,
	}, nil // todo (@pedro) - add validation
}

func (n *ContractDeployer) Start() error {
	fmt.Printf("Starting L1 contract deployer with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{"npx"}
	var ports []int

	// inspect stops operation until debugger is hooked on port 9229 if debug is enabled
	if n.cfg.DebugNamespaceEnabled {
		cmds = append(cmds, "--node-options=\"--inspect-brk=0.0.0.0:9229\"")
		ports = append(ports, 9229)
	}

	cmds = append(cmds, "hardhat", "deploy", "--network", "layer1")

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`
	{ 
        "layer1" : {
            "url" : "%s",
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ 
                "deployment_scripts/core"
            ],
            "accounts": [ "%s" ]
        }
    }
`, n.cfg.L1HTTPurl, n.cfg.PrivateKey),
	}

	containerID, err := docker.StartNewContainer(n.cfg.ContainerName, n.cfg.L1DeployerImage, cmds, ports, envs, nil, nil)
	if err != nil {
		return err
	}
	n.containerID = containerID
	return nil
}

func (n *ContractDeployer) RetrieveL1ContractAddresses() (config.NetworkConfig, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer func(cli *client.Client) {
		err := cli.Close()
		if err != nil {
			panic(err)
		}
	}(cli)

	// make sure the container has finished execution (3 minutes allows time for L1 transactions to be mined)
	err = docker.WaitForContainerToFinish(n.containerID, 3*time.Minute)
	if err != nil {
		return nil, err
	}

	tailSize := "3"
	if n.cfg.DebugNamespaceEnabled {
		tailSize = "4"
	}

	logsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tailSize,
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), n.containerID, logsOptions)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Buffer the output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		return nil, err
	}

	// Get the last lines
	output := buf.String()
	fmt.Printf("L2 Deployer output %s\n", output)

	lines := strings.Split(output, "\n")

	if n.cfg.DebugNamespaceEnabled {
		// remove debugger lines
		lines = lines[:len(lines)-2]
	}

	managementAddr, err := findAddress(lines[0])
	if err != nil {
		return nil, err
	}
	messageBusAddr, err := findAddress(lines[1])
	if err != nil {
		return nil, err
	}
	l1BlockHash := readValue("L1Start", lines[2])

	network := &config.NetworkInputConfig{
		ManagementContractAddress: managementAddr,
		MessageBusAddress:         messageBusAddr,
		L1StartHash:               l1BlockHash,
	}

	// Builds up the network configuration for later use depending on how it was assembled
	if n.netCfg != nil && n.netCfg.L1ChainID != "" {
		network.L1ChainID = n.netCfg.L1ChainID
	} else {
		network.L1ChainID = "1337"
		n.netCfg = network
	}

	return network, nil
}

func findAddress(line string) (string, error) {
	// Regular expression to match Ethereum addresses
	re := regexp.MustCompile("(0x[a-fA-F0-9]{40})")

	// Find all Ethereum addresses in the text
	matches := re.FindAllString(line, -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("no address found in: %s", line)
	}
	// Print the last
	return matches[len(matches)-1], nil
}

func readValue(name string, line string) string {
	parts := strings.Split(line, fmt.Sprintf("%s=", name))
	val := strings.TrimSpace(parts[len(parts)-1])
	return val
}
