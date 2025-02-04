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
	"github.com/sanity-io/litter"
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
	fmt.Printf("Starting L2 contract deployer with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{"/bin/sh"}
	var ports []int

	// inspect stops operation until debugger is hooked on port 9229 if debug is enabled
	if n.cfg.debugEnabled {
		cmds = append(cmds, "--node-options=\"--inspect-brk=0.0.0.0:9229\"")
		ports = append(ports, 9229)
	}

	cmds = append(cmds, "/home/obscuro/go-obscuro/entrypoint.sh", "obscuro:deploy", "--network", "layer2")

	envs := map[string]string{
		"L2_HOST":               n.cfg.l2Host,
		"L2_PORT":               strconv.Itoa(n.cfg.l2Port),
		"PREFUND_FAUCET_AMOUNT": n.cfg.faucetPrefundAmount,
		"MGMT_CONTRACT_ADDRESS": n.cfg.managementContractAddress,
		"MESSAGE_BUS_ADDRESS":   n.cfg.messageBusAddress,
		"NETWORK_JSON": fmt.Sprintf(`
{
        "layer1" : {
            "url" : "%s",
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ 
                "deployment_scripts/core"
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
				"InitialFunding",    	    "deployment_scripts/funding/layer1"
                "FaucetFunding",      		"deployment_scripts/funding/layer1"
                "CrossChainMessengerL1",    "deployment_scripts/messenger/layer1"
                "CrossChainMessengerL2",    "deployment_scripts/messenger/layer2
                "EthereumBridge",           "deployment_scripts/bridge/"
                "BridgeAdmin",              "deployment_scripts/testnet/layer1/"
                "ZenBase",                  "deployment_scripts/testnet/layer2/
                "SetFees"           	
            ],
            "accounts": [ 
                "%s"
            ]
        }
    }
`, n.cfg.l1HTTPURL, n.cfg.l1privateKey, n.cfg.l2PrivateKey),
	}

	containerID, err := docker.StartNewContainer("hh-l2-deployer", n.cfg.dockerImage, cmds, ports, envs, nil, nil, false)
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
