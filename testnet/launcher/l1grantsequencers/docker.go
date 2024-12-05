package l1grantsequencers

import (
	"fmt"
	"strings"

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
		"run",
		"--network",
		"layer1",
		"scripts/sequencer/001_grant_sequencers.ts",
		s.cfg.mgmtContractAddress,
		enclaveIDs,
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
	}

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
