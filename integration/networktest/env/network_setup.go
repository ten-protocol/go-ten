package env

import (
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

func Testnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://testnet.obscu.ro:80",
		[]string{"http://testnet.obscu.ro:80"}, // for now we'll just use sequencer as validator node... todo (@matt)
		"http://testnet-faucet.uksouth.azurecontainer.io/fund/obx",
		"ws://testnet-eth2network.uksouth.cloudapp.azure.com:9000",
	)
	return &testnetEnv{connector}
}

func DevTestnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://dev-testnet.obscu.ro:80",
		[]string{"http://dev-testnet.obscu.ro:80"}, // for now we'll just use sequencer as validator node... todo (@matt)
		"http://dev-testnet-faucet.uksouth.azurecontainer.io/fund/obx",
		"ws://dev-testnet-eth2network.uksouth.cloudapp.azure.com:9000",
	)
	return &testnetEnv{connector}
}

func LongRunningLocalNetwork(l1RPCAddress string) networktest.Environment {
	connector := NewTestnetConnectorWithFaucetAccount(
		"http://127.0.0.1:37800",
		[]string{"http://127.0.0.1:37801", "http://127.0.0.1:37802"},
		genesis.TestnetPrefundedPK,
		l1RPCAddress,
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
