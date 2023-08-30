package helpful

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

const (
	_sepoliaChainID = 11155111

	// To run Sepolia network: update these details with a websocket RPC address and funded PKs
	_sepoliaRPCAddress   = "wss://sepolia.infura.io/ws/v3/<api-key>"
	_sepoliaSequencerPK  = "<pk>" // account 0x<acc>
	_sepoliaValidator1PK = "<pk>" // account 0x<acc>

)

func TestRunLocalNetwork(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.EnsureTestLogsSetUp("local-geth-network")
	networkConnector, cleanUp, err := env.LocalDevNetwork().Prepare()
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
	checkBalance("sequencer", l1DeployerWallet, _sepoliaRPCAddress)

	val1Wallet := wallet.NewInMemoryWalletFromConfig(_sepoliaValidator1PK, _sepoliaChainID, testlog.Logger())
	checkBalance("validator1", val1Wallet, _sepoliaRPCAddress)

	validatorWallets := []wallet.Wallet{val1Wallet}
	networktest.EnsureTestLogsSetUp("local-network-live-l1")
	networkConnector, cleanUp, err := env.LocalNetworkLiveL1(l1DeployerWallet, validatorWallets, _sepoliaRPCAddress).Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUp()

	keepRunning(networkConnector)
}

func checkBalance(walDesc string, wal wallet.Wallet, rpcAddress string) {
	client, err := ethadapter.NewEthClientFromAddress(rpcAddress, 20*time.Second, common.HexToAddress("0x0"), testlog.Logger())
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
	fmt.Println("----")
	fmt.Println("Sequencer RPC", networkConnector.SequencerRPCAddress())
	for i := 0; i < networkConnector.NumValidators(); i++ {
		fmt.Println("Validator  ", i, networkConnector.ValidatorRPCAddress(i))
	}
	fmt.Println("----")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Network running until test is stopped...")
	<-done // Will block here until user hits ctrl+c
}

func TestReceiptSerialization(t *testing.T) {
	rec := types.Receipt{
		Type:              1,
		PostState:         []byte{1},
		Status:            1,
		CumulativeGasUsed: 123,
		Bloom:             types.Bloom{123},
		Logs:              nil,
		TxHash:            common.Hash{123},
		ContractAddress:   common.Address{123},
		GasUsed:           123,
		EffectiveGasPrice: big.NewInt(123),
		BlobGasUsed:       123,
		BlobGasPrice:      big.NewInt(123),
		BlockHash:         common.Hash{123},
		BlockNumber:       big.NewInt(123),
		TransactionIndex:  123,
	}

	recs := []*types.Receipt{&rec}
	ser, err := rlp.EncodeToBytes(recs)
	if err != nil {
		t.Fatal(err)
	}

	receipts := make(types.Receipts, 0)

	err = rlp.DecodeBytes(ser, &receipts)
	if err != nil {
		t.Fatal(err)
	}

	if len(receipts) != 1 {
		t.Fatal("unexpected number of receipts")
	}

	if receipts[0].Type != rec.Type {
		t.Fatal("unexpected receipt type")
	}

	if receipts[0].Status != rec.Status {
		fmt.Println(receipts[0].Status, rec.Status)
		t.Fatal("unexpected receipt status")
	}

	if receipts[0].ContractAddress != rec.ContractAddress {
		fmt.Println(receipts[0].ContractAddress, rec.ContractAddress)
		t.Fatal("unexpected receipt contract address")
	}

	if receipts[0].TxHash != rec.TxHash {
		fmt.Println(receipts[0].TxHash, rec.TxHash)
		t.Fatal("unexpected receipt tx hash")
	}

	if receipts[0].TransactionIndex != rec.TransactionIndex {
		fmt.Println(receipts[0].TransactionIndex, rec.TransactionIndex)
		t.Fatal("unexpected receipt transaction index")
	}
}
