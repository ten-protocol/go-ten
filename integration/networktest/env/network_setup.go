package env

import (
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

func Testnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://testnet.obscu.ro:13000",
		[]string{"http://testnet.obscu.ro:13000"}, // for now we'll just use sequencer as validator node... (todo)
		"http://testnet-faucet.uksouth.azurecontainer.io/fund/obx",
		"testnet-eth2network.uksouth.azurecontainer.io",
		9000,
	)
	return &testnetEnv{connector}
}

func DevTestnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://dev-testnet.obscu.ro:13000",
		[]string{"http://dev-testnet.obscu.ro:13000"}, // for now we'll just use sequencer as validator node... (todo)
		"http://dev-testnet-faucet.uksouth.azurecontainer.io/fund/obx",
		"dev-testnet-eth2network.uksouth.azurecontainer.io",
		9000,
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
