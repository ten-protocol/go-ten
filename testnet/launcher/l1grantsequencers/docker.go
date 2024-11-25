package l1grantsequencers

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/docker"
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
	cmds := []string{
		"npx",
		"run",
		"--network",
		"layer1",
		"scripts/sequencer/001_grant_sequencers.ts",
		s.cfg.mgmtContractAddress,
		s.cfg.enclaveIDs,
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
