package eth2network

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/common/docker"
	"github.com/obscuronet/go-obscuro/go/common/retry"
)

type Eth2Network struct {
	cfg *Config
}

func NewDockerEth2Network(cfg *Config) (*Eth2Network, error) {
	return &Eth2Network{
		cfg: cfg,
	}, nil // todo: add validation
}

func (n *Eth2Network) Start() error {
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

	_, err := docker.StartNewContainer("eth2network", "testnetobscuronet.azurecr.io/obscuronet/eth2network:latest", cmds, exposedPorts, nil, nil)
	return err
}

func (n *Eth2Network) IsReady() error {
	timeout := 10 * time.Minute
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
