package env

import (
	"fmt"

	"github.com/ten-protocol/go-ten/integration/common"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/tools/walletextension"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	// these ports were picked arbitrarily, if we want plan to use these tests on CI we need to use ports in the constants.go file
	_gwHTTPPort = 11180
	_gwWSPort   = 11181
)

func SepoliaTestnet(opts ...TestnetEnvOption) networktest.Environment {
	connector := newTestnetConnector(
		"http://erpc.sepolia-testnet.ten.xyz:80", // this is actually a validator...
		[]string{"http://erpc.sepolia-testnet.ten.xyz:80"},
		"http://sepolia-testnet-faucet.uksouth.azurecontainer.io/fund/eth",
		"wss://ethereum-sepolia-rpc.publicnode.com",
		"https://rpc.testnet.ten.xyz",  //"https://rpc.dexynth-gateway.ten.xyz",
		"wss://rpc.testnet.ten.xyz:81", // "wss://rpc.dexynth-gateway.ten.xyz:81",
	)
	return newTestnetEnv(connector, opts...)
}

func UATTestnet(opts ...TestnetEnvOption) networktest.Environment {
	connector := newTestnetConnector(
		"http://uat-sequencer.ten.xyz:80",
		[]string{"http://uat-validator-01.ten.xyz:80"},
		"http://uat-testnet-faucet.uksouth.azurecontainer.io/fund/eth",
		"wss://ethereum-sepolia-rpc.publicnode.com",
		"https://rpc.uat-testnet.ten.xyz",
		"wss://rpc.uat-testnet.ten.xyz:81",
	)
	return newTestnetEnv(connector, opts...)
}

func DevTestnet(opts ...TestnetEnvOption) networktest.Environment {
	connector := newTestnetConnector(
		"http://erpc.dev-testnet.ten.xyz:80", // this is actually a validator...
		[]string{"http://erpc.dev-testnet.ten.xyz:80"},
		"http://dev-testnet-faucet.uksouth.azurecontainer.io/fund/eth",
		"ws://dev-testnet-eth2network.uksouth.cloudapp.azure.com:9000",
		"https://rpc.dev-testnet.ten.xyz",
		"wss://rpc.dev-testnet.ten.xyz:81",
	)
	return newTestnetEnv(connector, opts...)
}

// LongRunningLocalNetwork is a local network, the l1WSURL is optional (can be empty string), only required if testing L1 interactions
func LongRunningLocalNetwork(l1WSURL string) networktest.Environment {
	connector := newTestnetConnectorWithFaucetAccount(
		"ws://127.0.0.1:17900",
		[]string{"ws://127.0.0.1:17901"},
		common.TestnetPrefundedPK,
		l1WSURL,
		"",
	)
	return newTestnetEnv(connector)
}

type TestnetEnvOption func(env *testnetEnv)

type testnetEnv struct {
	testnetConnector    *testnetConnector
	localTenGateway     bool
	tenGatewayContainer *walletextension.Container
	logger              gethlog.Logger
}

func (t *testnetEnv) Prepare() (networktest.NetworkConnector, func(), error) {
	if t.logger == nil {
		t.logger = testlog.Logger()
	}
	if t.localTenGateway {
		t.startTenGateway()
	}
	cleanup := func() {
		if t.tenGatewayContainer != nil {
			go func() {
				err := t.tenGatewayContainer.Stop()
				if err != nil {
					fmt.Println("failed to stop TEN gateway", err.Error())
				}
			}()
		}
	}
	// no cleanup or setup required for the testnet connector (unlike dev network which has teardown and startup to handle)
	return t.testnetConnector, cleanup, nil
}

func (t *testnetEnv) startTenGateway() {
	validator := t.testnetConnector.ValidatorRPCAddress(0)
	// remove http:// prefix for the gateway config
	validatorHTTP := validator[len("http://"):]
	// replace the last character with a 1 (expect it to be zero), this is good enough for these tests
	validatorWS := validatorHTTP[:len(validatorHTTP)-1] + "1"
	cfg := wecommon.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: _gwHTTPPort,
		WalletExtensionPortWS:   _gwWSPort,
		NodeRPCHTTPAddress:      validatorHTTP,
		NodeRPCWebsocketAddress: validatorWS,
		LogPath:                 "sys_out",
		LogLevel:                3, // info level
		DBType:                  "sqlite",
		TenChainID:              integration.TenChainID,
	}
	tenGWContainer := walletextension.NewContainerFromConfig(cfg, t.logger)

	fmt.Println("Starting TEN Gateway, HTTP Port:", _gwHTTPPort, "WS Port:", _gwWSPort)
	err := tenGWContainer.Start()
	if err != nil {
		t.logger.Error("failed to start TEN gateway", "err", err)
		panic(err)
	}
	t.tenGatewayContainer = tenGWContainer
	t.testnetConnector.tenGatewayURL = fmt.Sprintf("http://localhost:%d", _gwHTTPPort)
}

func newTestnetEnv(testnetConnector *testnetConnector, opts ...TestnetEnvOption) networktest.Environment {
	env := &testnetEnv{testnetConnector: testnetConnector}
	for _, opt := range opts {
		opt(env)
	}
	return env
}

func WithLocalTenGateway() TestnetEnvOption {
	return func(env *testnetEnv) {
		env.localTenGateway = true
	}
}
