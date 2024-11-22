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
	cmds := []string{"npx", "hardhat", "deploy", "--network", "layer1"}

	envs := map[string]string{
		"MGMT_CONTRACT_ADDRESS": s.cfg.mgmtContractAddress,
		"ENCLAVE_IDS":           s.cfg.enclaveIDs,
		"NETWORK_JSON": fmt.Sprintf(`{ 
            "layer1": {
                "url": "%s",
                "live": false,
                "saveDeployments": true,
                "deploy": [ 
                    "deployment_scripts/testnet/sequencer/"
                ],
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
