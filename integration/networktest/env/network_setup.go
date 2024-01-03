package env

import (
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

func SepoliaTestnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://erpc.sepolia-testnet.ten.xyz:80", // this is actually a validator...
		[]string{"http://erpc.sepolia-testnet.ten.xyz:80"},
		"http://sepolia-testnet-faucet.uksouth.azurecontainer.io/fund/eth",
		"https://rpc.sepolia.org/",
	)
	return &testnetEnv{connector}
}

func UATTestnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://erpc.uat-testnet.ten.xyz:80", // this is actually a validator...
		[]string{"http://erpc.uat-testnet.ten.xyz:80"},
		"http://uat-testnet-faucet.uksouth.azurecontainer.io/fund/eth",
		"ws://uat-testnet-eth2network.uksouth.cloudapp.azure.com:9000",
	)
	return &testnetEnv{connector}
}

func DevTestnet() networktest.Environment {
	connector := NewTestnetConnector(
		"http://erpc.dev-testnet.ten.xyz:80", // this is actually a validator...
		[]string{"http://erpc.dev-testnet.ten.xyz:80"},
		"http://dev-testnet-faucet.uksouth.azurecontainer.io/fund/eth",
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
