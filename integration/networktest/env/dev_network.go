package env

import (
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/simulation/devnetwork"
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

func awaitNodesAvailable(_ networktest.NetworkConnector) error { //nolint:unparam
	// todo: create RPC clients for all the nodes and wait until their health checks pass

	// for now we just sleep
	time.Sleep(15 * time.Second)
	return nil
}

func LocalDevNetwork() networktest.Environment {
	return &devNetworkEnv{}
}
