package l1contractdeployer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/ten-protocol/go-ten/go/node"

	"github.com/docker/docker/client"
	"github.com/ten-protocol/go-ten/go/common/docker"
)

type ContractDeployer struct {
	cfg         *Config
	containerID string
}

func NewDockerContractDeployer(cfg *Config) (*ContractDeployer, error) {
	return &ContractDeployer{
		cfg: cfg,
	}, nil // todo (@pedro) - add validation
}

func (n *ContractDeployer) Start() error {
	fmt.Printf("Starting L1 contract deployer with config: \n%v\n\n", n.cfg)

	cmds := []string{"npx"}
	var ports []int

	// inspect stops operation until debugger is hooked on port 9229 if debug is enabled
	if n.cfg.DebugEnabled {
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
                "deployment_scripts/core",
				"deployment_scripts/testnet/layer1"
            ],
            "accounts": [ "%s" ]
        }
    }
`, n.cfg.L1HTTPURL, n.cfg.PrivateKey),
	}

	containerID, err := docker.StartNewContainer("hh-l1-deployer", n.cfg.DockerImage, cmds, ports, envs, nil, nil, false)
	if err != nil {
		return err
	}
	n.containerID = containerID
	return nil
}

func (n *ContractDeployer) RetrieveL1ContractAddresses() (*node.NetworkConfig, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// make sure the container has finished execution (5 minutes allows time for L1 transactions to be mined)
	err = docker.WaitForContainerToFinish(n.containerID, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	tailSize := "7"
	if n.cfg.DebugEnabled {
		tailSize = "8"
	}

	logsOptions := container.LogsOptions{
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

	if n.cfg.DebugEnabled {
		// remove debugger lines
		lines = lines[:len(lines)-2]
	}

	networkConfigAddr, err := findAddress(lines[0])
	if err != nil {
		return nil, err
	}
	crossChainAddr, err := findAddress(lines[1])
	if err != nil {
		return nil, err
	}
	messageBusAddr, err := findAddress(lines[2])
	if err != nil {
		return nil, err
	}
	enclaveRegistryAddr, err := findAddress(lines[3])
	if err != nil {
		return nil, err
	}
	daRegistryAddr, err := findAddress(lines[4])
	if err != nil {
		return nil, err
	}
	bridgeAddress, err := findAddress(lines[5])
	if err != nil {
		return nil, err
	}
	l1BlockHash := readValue("L1Start", lines[6])

	return &node.NetworkConfig{
		EnclaveRegistryAddress:          enclaveRegistryAddr,
		DataAvailabilityRegistryAddress: daRegistryAddr,
		CrossChainAddress:               crossChainAddr,
		NetworkConfigAddress:            networkConfigAddr,
		MessageBusAddress:               messageBusAddr,
		L1StartHash:                     l1BlockHash,
		BridgeAddress:                   bridgeAddress,
	}, nil
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
