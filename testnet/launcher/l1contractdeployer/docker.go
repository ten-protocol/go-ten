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
		"SEQUENCER_HOST_ADDRESS": n.cfg.SequencerHostAddress,
		"ETHERSCAN_API_KEY":      n.cfg.EtherscanAPIKey,
		"MAX_GAS_GWEI":           n.cfg.MaxGasGwei,
		"CHECK_GAS_PRICE":        n.cfg.CheckGasPrice,
		"USDC_ADDRESS":           n.cfg.USDCAddress,
		"USDT_ADDRESS":           n.cfg.USDTAddress,
		"WETH_ADDRESS":           n.cfg.WETHAddress,
		"NETWORK_JSON": fmt.Sprintf(`
{
        "layer1" : {
            "url" : "%s",
            "useGateway" : false,
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

	// make sure the container has finished execution
	// (generous 30 minute timeout allows time for L1 transactions to be mined in unpredictable environments)
	err = docker.WaitForContainerToFinish(n.containerID, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	logsOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "50", // fetch more lines to ensure we capture all contract addresses
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

	output := buf.String()
	fmt.Printf("L1 Deployer output:\n%s\n", output)

	// Parse contract addresses by searching for specific keys in the output
	// This is more robust than relying on exact line positions
	networkConfigAddr, err := findAddressByKey(output, "NetworkConfig")
	if err != nil {
		return nil, fmt.Errorf("failed to find NetworkConfig address: %w", err)
	}
	crossChainAddr, err := findAddressByKey(output, "CrossChain")
	if err != nil {
		return nil, fmt.Errorf("failed to find CrossChain address: %w", err)
	}
	messageBusAddr, err := findAddressByKey(output, "MerkleMessageBus")
	if err != nil {
		return nil, fmt.Errorf("failed to find MerkleMessageBus address: %w", err)
	}
	enclaveRegistryAddr, err := findAddressByKey(output, "NetworkEnclaveRegistry")
	if err != nil {
		return nil, fmt.Errorf("failed to find NetworkEnclaveRegistry address: %w", err)
	}
	daRegistryAddr, err := findAddressByKey(output, "DataAvailabilityRegistry")
	if err != nil {
		return nil, fmt.Errorf("failed to find DataAvailabilityRegistry address: %w", err)
	}
	bridgeAddress, err := findAddressByKey(output, "L1Bridge")
	if err != nil {
		return nil, fmt.Errorf("failed to find L1Bridge address: %w", err)
	}
	l1BlockHash, err := findValueByKey(output, "L1Start")
	if err != nil {
		return nil, fmt.Errorf("failed to find L1Start hash: %w", err)
	}

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

// findAddressByKey searches for a line containing "key=" and extracts the Ethereum address from it.
// This is more robust than relying on exact line positions in the output.
func findAddressByKey(output, key string) (string, error) {
	// Look for pattern like "NetworkConfig= 0x..." or "NetworkConfig=0x..."
	keyPattern := regexp.MustCompile(key + `=\s*(0x[a-fA-F0-9]{40})`)
	matches := keyPattern.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// Fallback: find the line containing the key and extract any address from it
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, key+"=") {
			addrPattern := regexp.MustCompile(`(0x[a-fA-F0-9]{40})`)
			addrMatches := addrPattern.FindStringSubmatch(line)
			if len(addrMatches) >= 1 {
				return addrMatches[1], nil
			}
		}
	}

	return "", fmt.Errorf("no address found for key %s in output", key)
}

// findValueByKey searches for a line containing "key=" and extracts the value after it.
// Used for non-address values like block hashes.
func findValueByKey(output, key string) (string, error) {
	// Look for pattern like "L1Start= 0x..." or "L1Start=0x..."
	keyPattern := regexp.MustCompile(key + `=\s*(0x[a-fA-F0-9]+)`)
	matches := keyPattern.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// Fallback: find the line containing the key and extract the value
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, key+"=") {
			parts := strings.Split(line, key+"=")
			if len(parts) >= 2 {
				val := strings.TrimSpace(parts[1])
				// Extract just the hex value (stop at whitespace or newline)
				hexPattern := regexp.MustCompile(`(0x[a-fA-F0-9]+)`)
				hexMatches := hexPattern.FindStringSubmatch(val)
				if len(hexMatches) >= 1 {
					return hexMatches[1], nil
				}
			}
		}
	}

	return "", fmt.Errorf("no value found for key %s in output", key)
}
