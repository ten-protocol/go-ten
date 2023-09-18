package env

import (
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

func Testnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://erpc.testnet.obscu.ro:80",
		[]string{"http://erpc.testnet.obscu.ro:80"}, // for now we'll just use sequencer as validator node... todo (@matt)
		"http://testnet-faucet.uksouth.azurecontainer.io/fund/obx",
		"ws://testnet-eth2network.uksouth.cloudapp.azure.com:9000",
	)
	return &testnetEnv{connector}
}

func DevTestnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://erpc.dev-testnet.obscu.ro:80",
		[]string{"http://erpc.dev-testnet.obscu.ro:80"}, // for now we'll just use sequencer as validator node... todo (@matt)
		"http://dev-testnet-faucet.uksouth.azurecontainer.io/fund/obx",
		"ws://dev-testnet-eth2network.uksouth.cloudapp.azure.com:9000",
	)
	return &testnetEnv{connector}
}

// LongRunningLocalNetwork is a local network, the l1WSURL is optional (can be empty string), only required if testing L1 interactions
func LongRunningLocalNetwork(l1WSURL string) networktest.Environment {
	connector := NewTestnetConnectorWithFaucetAccount(
		"ws://127.0.0.1:37900",
		[]string{"ws://127.0.0.1:37901"},
		genesis.TestnetPrefundedPK,
		l1WSURL,
	)
	return &testnetEnv{connector}
}

type testnetEnv struct {
	testnetConnector networktest.NetworkConnector
}

func (t *testnetEnv) Prepare() (networktest.NetworkConnector, func(), error) {
	// no cleanup or setup required for the testnet connector (unlike dev network which has teardown and startup to handle)
	return t.testnetConnector, func() {}, nil
}
