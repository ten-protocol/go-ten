package l2contractdeployer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"

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
	}, nil
}

func (n *ContractDeployer) Start() error {
	cmds := []string{"/bin/sh"}
	var ports []int

	// inspect stops operation until debugger is hooked on port 9229 if debug is enabled
	if n.cfg.DebugEnabled {
		cmds = append(cmds, "--node-options=\"--inspect-brk=0.0.0.0:9229\"")
		ports = append(ports, 9229)
	}

	cmds = append(cmds, "/home/obscuro/go-obscuro/entrypoint.sh", "obscuro:deploy", "--network", "layer2")

	envs := map[string]string{
		"L2_HOST":               n.cfg.L2Host,
		"L2_HTTP_PORT":          strconv.Itoa(n.cfg.L2HTTPPort),
		"L2_WS_PORT":            strconv.Itoa(n.cfg.L2WSPort),
		"PREFUND_FAUCET_AMOUNT": n.cfg.FaucetPrefundAmount,
		"ENCLAVE_REGISTRY_ADDR": n.cfg.EnclaveRegistryAddress,
		"CROSS_CHAIN_ADDR":      n.cfg.CrossChainAddress,
		"DA_REGISTRY_ADDR":      n.cfg.DaRegistryAddress,
		"NETWORK_CONFIG_ADDR":   n.cfg.NetworkConfigAddress,
		"MESSAGE_BUS_ADDR":      n.cfg.MessageBusAddress,
		"NETWORK_CHAINID":       strconv.FormatInt(n.cfg.ChainID, 10),
		"NETWORK_JSON": fmt.Sprintf(`
{
        "layer1" : {
            "url" : "%s",
            "gasMultiplier" : 1.2,
            "useGateway" : false,
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ 
                "deployment_scripts/core",
                "deployment_scripts/testnet/layer1"
            ],
            "accounts": [ 
                "%s"
            ]
        },
        "layer2" : {
            "url": "http://127.0.0.1:3000/v1/",
			"useGateway": true,
            "live" : false,
            "saveDeployments" : true,
            "companionNetworks" : { "layer1" : "layer1" },
            "deploy": [ 
				"deployment_scripts/funding/layer1",
                "deployment_scripts/bridge/",
                "deployment_scripts/testnet/layer2/"
            ],
            "accounts": [ 
                "%s"
            ]
        }
    }
`, n.cfg.L1HTTPURL, n.cfg.L1PrivateKey, n.cfg.L2PrivateKey),
	}

	containerID, err := docker.StartNewContainer("hh-l2-deployer", n.cfg.DockerImage, cmds, ports, envs, nil, nil, false)
	if err != nil {
		return err
	}
	n.containerID = containerID
	return nil
}

func (n *ContractDeployer) GetID() string {
	return n.containerID
}

func (n *ContractDeployer) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	// make sure the container has finished execution
	err = docker.WaitForContainerToFinish(n.containerID, 15*time.Minute)
	if err != nil {
		n.PrintLogs(cli)
		return err
	}

	// if we want to read anything from the container logs we can do it here (see RetrieveL1ContractAddresses as example)

	return nil
}

func (n *ContractDeployer) PrintLogs(cli *client.Client) {
	logsOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), n.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", n.containerID, err)
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", n.containerID)
	}
	fmt.Println(buf.String())
}
