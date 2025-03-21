package l1grantsequencers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
)

type GrantSequencers struct {
	cfg         *Config
	containerID string
}

func NewGrantSequencers(cfg *Config) (*GrantSequencers, error) {
	return &GrantSequencers{
		cfg: cfg,
	}, nil
}

func (s *GrantSequencers) Start() error {
	fmt.Printf("Starting grant sequencers with config: %s\n", s.cfg)
	var enclaveIDs string
	var err error
	if s.cfg.enclaveIDs != "" {
		enclaveIDs = s.cfg.enclaveIDs
	} else if s.cfg.sequencerURL != "" {
		enclaveIDs, err = fetchEnclaveIDs(s.cfg.sequencerURL)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("enclaveIDs or sequencerURL must be provided")
	}
	cmds := []string{
		"npx",
		"hardhat",
		"run",
		"--network",
		"layer1",
		"scripts/sequencer/001_grant_sequencers.ts",
	}

	envs := map[string]string{
		"NETWORK_JSON": fmt.Sprintf(`{ 
            "layer1": {
                "url": "%s",
                "live": false,
                "saveDeployments": true,
                "accounts": [ "%s" ]
            }
        }`, s.cfg.l1HTTPURL, s.cfg.privateKey),
		"ENCLAVE_REGISTRY_ADDR": s.cfg.enclaveRegistryAddress,
		"ENCLAVE_IDS":           enclaveIDs,
	}

	fmt.Printf("Starting grant sequencer script. EnclaveRegistryAddress: %s, EnclaveIDs: %s\n", s.cfg.enclaveRegistryAddress, enclaveIDs)

	containerID, err := docker.StartNewContainer(
		"grant-sequencers",
		s.cfg.dockerImage,
		cmds,
		nil,
		envs,
		nil,
		nil,
		false,
	)
	if err != nil {
		return err
	}
	s.containerID = containerID
	return nil
}

func (s *GrantSequencers) WaitForFinish() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	// make sure the container has finished execution
	err = docker.WaitForContainerToFinish(s.containerID, 15*time.Minute)
	if err != nil {
		fmt.Println("Error waiting for container to finish: ", err)
		s.PrintLogs(cli)
		return err
	}

	return nil
}

func fetchEnclaveIDs(url string) (string, error) {
	// fetch enclaveIDs
	client, err := rpc.NewNetworkClient(url)
	if err != nil {
		return "", fmt.Errorf("failed to create network client (%s): %w", url, err)
	}
	defer client.Stop()

	obsClient := obsclient.NewObsClient(client)
	health, err := obsClient.Health()
	if err != nil {
		return "", fmt.Errorf("failed to get health status: %w", err)
	}

	if len(health.Enclaves) == 0 {
		return "", fmt.Errorf("could not retrieve enclave IDs from health endpoint - no enclaves found")
	}

	var enclaveIDs []string
	for _, status := range health.Enclaves {
		enclaveIDs = append(enclaveIDs, status.EnclaveID.String())
	}
	return strings.Join(enclaveIDs, ","), nil
}

func (s *GrantSequencers) PrintLogs(cli *client.Client) {
	logsOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	// Read the container logs
	out, err := cli.ContainerLogs(context.Background(), s.containerID, logsOptions)
	if err != nil {
		fmt.Printf("Error printing out container %s logs... %v\n", s.containerID, err)
	}
	defer out.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil {
		fmt.Printf("Error getting logs for container %s\n", s.containerID)
	}
	fmt.Println(buf.String())
}
