package helpful

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

const (
	_sepoliaChainID = 11155111

	SepoliaRPCAddress1 = "wss://sepolia.infura.io/ws/v3/<api-key>" // seq
	SepoliaRPCAddress2 = "wss://sepolia.infura.io/ws/v3/<api-key>" // val
	SepoliaRPCAddress3 = "wss://sepolia.infura.io/ws/v3/<api-key>" // tester

	_sepoliaSequencerPK  = "<pk>" // account 0x<acc>
	_sepoliaValidator1PK = "<pk>" // account 0x<acc>
)

func TestRunLocalNetwork(t *testing.T) {
	//	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-geth-network")
	networkConnector, cleanUp, err := env.LocalDevNetwork(devnetwork.WithGateway()).Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func TestRunLocalGatewayAgainstRemoteTestnet(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-faucet-remote-testnet")

	// set the testnet the gateway will connect to here
	netw := env.SepoliaTestnet(env.WithLocalTenGateway())
	networkConnector, cleanUp, err := netw.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func TestRunLocalNetworkAgainstSepolia(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-sepolia-network")

	l1DeployerWallet := wallet.NewInMemoryWalletFromConfig(_sepoliaSequencerPK, _sepoliaChainID, testlog.Logger())
	checkBalance("sequencer", l1DeployerWallet, SepoliaRPCAddress1)

	val1Wallet := wallet.NewInMemoryWalletFromConfig(_sepoliaValidator1PK, _sepoliaChainID, testlog.Logger())
	checkBalance("validator1", val1Wallet, SepoliaRPCAddress2)

	validatorWallets := []wallet.Wallet{val1Wallet}
	networktest.EnsureTestLogsSetUp("local-network-live-l1")
	networkConnector, cleanUp, err := env.LocalNetworkLiveL1(l1DeployerWallet, validatorWallets, []string{SepoliaRPCAddress1, SepoliaRPCAddress2}).Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func checkBalance(walDesc string, wal wallet.Wallet, rpcAddress string) {
	client, err := ethadapter.NewEthClientFromURL(rpcAddress, 20*time.Second, common.HexToAddress("0x0"), testlog.Logger())
	if err != nil {
		panic("unable to create live L1 eth client, err=" + err.Error())
	}

	bal, err := client.BalanceAt(wal.Address(), nil)
	if err != nil {
		panic(fmt.Errorf("failed to get balance for %s (%s): %w", walDesc, wal.Address(), err))
	}
	fmt.Println(walDesc, "wallet balance", wal.Address(), bal)

	if bal.Cmp(big.NewInt(0)) <= 0 {
		panic(fmt.Errorf("%s wallet has no funds: %s", walDesc, wal.Address()))
	}
}

func keepRunning(networkConnector networktest.NetworkConnector) {
	gatewayURL, err := networkConnector.GetGatewayURL()
	hasGateway := err == nil

	fmt.Println("----")
	fmt.Println("Sequencer RPC", networkConnector.SequencerRPCAddress())
	for i := 0; i < networkConnector.NumValidators(); i++ {
		fmt.Println("Validator  ", i, networkConnector.ValidatorRPCAddress(i))
	}
	if hasGateway {
		fmt.Println("Gateway      ", gatewayURL)
	}
	fmt.Println("----")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Network running until test is stopped...")
	<-done // Will block here until user hits ctrl+c
}
