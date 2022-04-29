//go:build ganache
// +build ganache

package simulation

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/contracts"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// TestMemObscuroRealEthMonteCarloSimulation runs the simulation against a ganache network
// 1. Install ganache -> npm install ganache --global
// 2. Run ganache -> rm -rf ganachedb &&  ganache --database.dbPath="./ganachedb"  -l 1024000000000 --wallet.accounts="0x5dbbff1b5ff19f1ad6ea656433be35f6846e890b3f3ec6ef2b2e2137a8cab4ae,0x56BC75E2D63100000" --wallet.accounts="0xb728cd9a9f54cede03a82fc189eab4830a612703d48b7ef43ceed2cbad1a06c7,0x56BC75E2D63100000" --wallet.accounts="0x1e1e76d5c0ea1382b6acf76e873977fd223c7fa2a6dc57db2b94e93eb303ba85,0x56BC75E2D63100000" -p 7545 -g 225 --miner.callGasLimit 1024000000000
func TestMemObscuroRealEthMonteCarloSimulation(t *testing.T) {
	setupTestLog()

	// private key is prefunded and used to issue txs - used here to deploy contract ahead of node initialization
	tmpWallet := wallet.NewInMemoryWallet("1e1e76d5c0ea1382b6acf76e873977fd223c7fa2a6dc57db2b94e93eb303ba85")
	contractAddr := deployContract(t, tmpWallet)

	params := params.SimParams{
		NumberOfNodes:             2,
		NumberOfWallets:           2,
		AvgBlockDuration:          time.Second,
		SimulationTime:            15 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.32,
		L2ToL1EfficiencyThreshold: 0.4,
		TxHandler:                 mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr),
		MgmtContractAddr:          contractAddr,
	}

	params.AvgNetworkLatency = params.AvgBlockDuration / 15
	params.AvgGossipPeriod = params.AvgBlockDuration / 2

	testSimulation(t, network.NewNetworkInMemoryGeth(), params)
}

func deployContract(t *testing.T, w wallet.Wallet) common.Address {
	tmpClient, err := ethclient.NewEthClient(common.Address{}, "127.0.0.1", 7545, w, common.Address{})
	if err != nil {
		t.Fatal(err)
	}

	deployContractTx := types.LegacyTx{
		Nonce:    0, // relies on a clean env
		GasPrice: big.NewInt(2000000000),
		Gas:      1025_000_000,
		Data:     common.Hex2Bytes(contracts.MgmtContractByteCode),
	}

	signedTx, err := tmpClient.SubmitTransaction(&deployContractTx)
	if err != nil {
		t.Fatal(err)
	}

	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < 10*time.Second; time.Sleep(time.Second) {
		receipt, err = tmpClient.FetchTxReceipt(signedTx.Hash())
		if err == nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			t.Fatal(err)
		}
		t.Logf("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	t.Logf("Contract deployed to %s\n", receipt.ContractAddress)
	return receipt.ContractAddress
}
