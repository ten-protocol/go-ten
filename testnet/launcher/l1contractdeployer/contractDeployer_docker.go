package l1contractdeployer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/obscuronet/go-obscuro/go/common/docker"
)

type ContractDeployer struct {
	cfg         *Config
	containerID string
}

func NewDockerContractDeployer(cfg *Config) (*ContractDeployer, error) {
	return &ContractDeployer{
		cfg: cfg,
	}, nil // todo: add validation
}

func (n *ContractDeployer) Start() error {
	cmds := []string{
		"npx", "hardhat", "deploy",
		"--network", "layer1",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`
{ 
        "layer1" : {
            "url" : "http://%s:%d",
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ 
                "deployment_scripts/core"
            ],
            "accounts": [ "%s" ]
        }
    }
`, n.cfg.l1Host, n.cfg.l1Port, n.cfg.privateKey),
	}

	containerID, err := docker.StartNewContainer("hh-l1-deployer", "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", cmds, nil, envs)
	if err != nil {
		return err
	}
	n.containerID = containerID
	return nil
}

func (n *ContractDeployer) RetrieveL1ContractAddresses() (string, string, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", "", err
	}
	defer cli.Close()

	// make sure the container has finished execution
	err = docker.WaitForContainerToFinish(n.containerID, time.Minute)
	if err != nil {
		return "", "", err
	}

	logsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "2",
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), n.containerID, logsOptions)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Buffer the output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		panic(err)
	}

	// Get the last two lines
	output := buf.String()
	lines := strings.Split(output, "\n")

	managementAddr, err := findAddress(lines[0])
	if err != nil {
		return "", "", err
	}
	messageBusAddr, err := findAddress(lines[1])
	if err != nil {
		return "", "", err
	}

	return managementAddr, messageBusAddr, nil
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
