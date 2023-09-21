package l2contractdeployer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/obscuronet/go-obscuro/go/common/docker"
	"github.com/sanity-io/litter"
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

	cmds := []string{
		"npx", "hardhat", "obscuro:deploy",
		"--network", "layer2",
	}

	envs := map[string]string{
		"MESSAGE_BUS_ADDRESS": n.cfg.messageBusAddress,
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
            "obscuroEncRpcUrl" : "ws://%s:%d",
            "url": "http://127.0.0.1:3000/v1",
            "live" : false,
            "saveDeployments" : true,
            "companionNetworks" : { "layer1" : "layer1" },
            "deploy": [ 
				"deployment_scripts/funding/layer1",
                "deployment_scripts/messenger/layer1",
                "deployment_scripts/messenger/layer2",
                "deployment_scripts/bridge/",
                "deployment_scripts/testnet/layer1/",
                "deployment_scripts/testnet/layer2/"
            ],
            "accounts": [ 
                "%s",
                "%s",
                "%s"
            ]
        }
    }
`, n.cfg.l1HTTPURL, n.cfg.l1privateKey, n.cfg.l2Host, n.cfg.l2Port, n.cfg.l2PrivateKey, n.cfg.hocPKString, n.cfg.pocPKString),
	}

	containerID, err := docker.StartNewContainer("hh-l2-deployer", n.cfg.dockerImage, cmds, nil, envs, nil, nil)
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
	err = docker.WaitForContainerToFinish(n.containerID, 3*time.Minute)
	if err != nil {
		n.PrintLogs(cli)
		return err
	}

	// if we want to read anything from the container logs we can do it here (see RetrieveL1ContractAddresses as example)

	return nil
}

func (n *ContractDeployer) PrintLogs(cli *client.Client) {
	logsOptions := types.ContainerLogsOptions{
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
