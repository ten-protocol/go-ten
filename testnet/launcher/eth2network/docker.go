package eth2network

import (
	"context"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sanity-io/litter"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/common/retry"
)

type Eth2Network struct {
	cfg *config.Eth2NetworkConfig
}

func NewDockerEth2Network(cfg *config.TestnetConfig) (*Eth2Network, error) {
	return &Eth2Network{
		cfg: &cfg.Eth2Network,
	}, nil // todo (@pedro) - add validation
}

func (n *Eth2Network) Start() error {
	fmt.Printf("Starting Eth2Network with config: \n%s\n\n", litter.Sdump(*n.cfg))

	cmds := []string{
		"/home/obscuro/go-obscuro/integration/eth2network/main/main",
		"--numNodes", strconv.Itoa(n.cfg.GethNumNodes),
	}

	if len(n.cfg.GethPrefundedAddresses) > 1 {
		cmds = append(cmds, "-prefundedAddrs", strings.Join(n.cfg.GethPrefundedAddresses, ","))
	}

	var exposedPorts []int
	if n.cfg.GethHTTPPort != 0 {
		cmds = append(cmds, "-gethHTTPStartPort", fmt.Sprintf("%d", n.cfg.GethHTTPPort))
		exposedPorts = append(exposedPorts, n.cfg.GethHTTPPort)
	}

	if n.cfg.GethWebsocketPort != 0 {
		cmds = append(cmds, "-gethWSStartPort", fmt.Sprintf("%d", n.cfg.GethWebsocketPort))
		exposedPorts = append(exposedPorts, n.cfg.GethWebsocketPort)
	}

	_, err := docker.StartNewContainer(n.cfg.ContainerName, n.cfg.GethImage, cmds, exposedPorts, nil, nil, nil)
	return err
}

func (n *Eth2Network) IsReady() error {
	timeout := 20 * time.Minute // this can be reduced when we no longer download the ethereum binaries
	interval := 2 * time.Second
	var dial *ethclient.Client
	var err error

	// retry the connection
	err = retry.Do(func() error {
		dial, err = ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.cfg.GethHTTPPort))
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
