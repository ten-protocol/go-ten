package eth2network

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/common/retry"
)

type Eth2Network struct {
	cfg *Config
}

func NewDockerEth2Network(cfg *Config) (*Eth2Network, error) {
	return &Eth2Network{
		cfg: cfg,
	}, nil // todo (@pedro) - add validation
}

func (n *Eth2Network) Start() error {
	fmt.Printf("Starting Eth2Network with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{
		"/home/obscuro/go-obscuro/integration/eth2network/main/main",
		"--numNodes", "1",
	}

	if len(n.cfg.prefundedAddrs) > 1 {
		cmds = append(cmds, "-prefundedAddrs", strings.Join(n.cfg.prefundedAddrs, ","))
	}

	var exposedPorts []int
	if n.cfg.gethHTTPPort != 0 {
		cmds = append(cmds, "-gethHTTPStartPort", fmt.Sprintf("%d", n.cfg.gethHTTPPort))
		exposedPorts = append(exposedPorts, n.cfg.gethHTTPPort)
	}

	if n.cfg.gethWSPort != 0 {
		cmds = append(cmds, "-gethWSStartPort", fmt.Sprintf("%d", n.cfg.gethWSPort))
		exposedPorts = append(exposedPorts, n.cfg.gethWSPort)
	}

	// keep a volume of binaries to avoid downloading
	volume := map[string]string{"eth2_bin": "/home/obscuro/go-obscuro/integration/.build/eth2_bin/"}

	_, err := docker.StartNewContainer("eth2network", "testnetobscuronet.azurecr.io/obscuronet/eth2network:latest", cmds, exposedPorts, nil, nil, volume)
	return err
}

func (n *Eth2Network) IsReady() error {
	timeout := 20 * time.Minute // this can be reduced when we no longer download the ethereum binaries
	interval := 2 * time.Second
	var dial *ethclient.Client
	var err error

	// retry the connection
	err = retry.Do(func() error {
		dial, err = ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.cfg.gethHTTPPort))
		if err != nil {
			return err
		}
		return nil
	}, retry.NewTimeoutStrategy(timeout, interval))

	// wait until merge block
	return retry.Do(func() error {
		number, err := dial.BlockNumber(context.Background())
		if err != nil {
			return err
		}

		if number <= 7 {
			return fmt.Errorf("retry - post-merge block has not been reached yet - current block: %d", number)
		}

		return nil
	}, retry.NewTimeoutStrategy(timeout, interval))
}
