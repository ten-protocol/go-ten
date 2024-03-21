package env

import (
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/wallet"

	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

type devNetworkEnv struct {
	inMemDevNetwork *devnetwork.InMemDevNetwork
}

func (d *devNetworkEnv) Prepare() (networktest.NetworkConnector, func(), error) {
	d.inMemDevNetwork.Start()

	err := awaitNodesAvailable(d.inMemDevNetwork)
	if err != nil {
		return nil, nil, err
	}

	return d.inMemDevNetwork, d.inMemDevNetwork.CleanUp, nil
}

func awaitNodesAvailable(nc networktest.NetworkConnector) error {
	err := awaitHealthStatus(nc.GetSequencerNode().HostRPCWSAddress(), 60*time.Second)
	if err != nil {
		return err
	}
	for i := 0; i < nc.NumValidators(); i++ {
		err := awaitHealthStatus(nc.GetValidatorNode(i).HostRPCWSAddress(), 60*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

// awaitHealthStatus waits for the host to be healthy until timeout
func awaitHealthStatus(rpcAddress string, timeout time.Duration) error {
	fmt.Println("Awaiting health status:", rpcAddress)
	return retry.Do(func() error {
		c, err := obsclient.Dial(rpcAddress)
		if err != nil {
			return fmt.Errorf("failed dial host (%s): %w", rpcAddress, err)
		}
		defer c.Close()
		healthy, err := c.Health()
		if err != nil {
			return fmt.Errorf("failed to get host health (%s): %w", rpcAddress, err)
		}
		if !healthy {
			return fmt.Errorf("host is not healthy (%s)", rpcAddress)
		}
		return nil
	}, retry.NewTimeoutStrategy(timeout, 200*time.Millisecond))
}

func LocalDevNetwork(opts ...devnetwork.TenConfigOption) networktest.Environment {
	return &devNetworkEnv{inMemDevNetwork: devnetwork.LocalDevNetwork(opts...)}
}

// LocalNetworkLiveL1 creates a local network that points to a live running L1.
// Note: seqWallet and validatorWallets need funds. seqWallet is used to deploy the L1 contracts
func LocalNetworkLiveL1(seqWallet wallet.Wallet, validatorWallets []wallet.Wallet, l1RPCURLs []string) networktest.Environment {
	return &devNetworkEnv{inMemDevNetwork: devnetwork.LiveL1DevNetwork(seqWallet, validatorWallets, l1RPCURLs)}
}
