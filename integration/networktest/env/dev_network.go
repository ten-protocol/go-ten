package env

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/simulation/devnetwork"
	"time"
)

type devNetworkEnv struct{}

func (d *devNetworkEnv) Prepare() (networktest.NetworkConnector, func(), error) {
	devNet := devnetwork.DefaultDevNetwork()
	devNet.Start()

	err := awaitNodesAvailable(devNet)
	if err != nil {
		return nil, nil, err
	}

	return devNet, devNet.CleanUp, nil
}

func awaitNodesAvailable(nc networktest.NetworkConnector) error { //nolint:unparam
	err := awaitHealthStatus(nc.GetSequencerNode().HostRPCAddress(), 30*time.Second)
	if err != nil {
		return err
	}
	for i := 0; i < nc.NumValidators(); i++ {
		err := awaitHealthStatus(nc.GetValidatorNode(i).HostRPCAddress(), 30*time.Second)
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

func LocalDevNetwork() networktest.Environment {
	return &devNetworkEnv{}
}
