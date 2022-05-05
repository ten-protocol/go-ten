package simulation

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// TestGethMemObscuroEthMonteCarloSimulation runs the simulation against a private geth network using Clique (PoA)
func TestGethMemObscuroEthMonteCarloSimulation(t *testing.T) {
	setupTestLog()

	numberOfNodes := 5

	// randomly create the ethereum wallets to be used and prefund them
	wallets := make([]wallet.Wallet, numberOfNodes)
	walletAddresses := make([]string, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		wallets[i] = datagenerator.RandomWallet()
		walletAddresses[i] = wallets[i].Address().String()
	}

	// make sure the network binaries exist
	path, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		t.Fatal(err)
	}

	// kickoff the network with the prefunded wallet addresses
	gethNetwork := gethnetwork.NewGethNetwork(
		40000,
		path,
		numberOfNodes,
		1,
		walletAddresses,
	)
	defer gethNetwork.StopNodes()

	// take the first random wallet and deploy the contract in the network
	contractAddr := deployContract(t, wallets[0], gethNetwork.WebSocketPorts[0])

	params := params.SimParams{
		NumberOfNodes:             numberOfNodes,
		NumberOfObscuroWallets:    10,
		AvgBlockDuration:          time.Second,
		SimulationTime:            25 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.5,
		L2ToL1EfficiencyThreshold: 0.9,
		TxHandler:                 mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr),
		MgmtContractAddr:          contractAddr,
		EthWallets:                wallets,
	}

	params.AvgNetworkLatency = params.AvgBlockDuration / 50
	params.AvgGossipPeriod = params.AvgBlockDuration / 3

	testSimulation(t, network.NewNetworkInMemoryGeth(&gethNetwork), params)
}

func deployContract(t *testing.T, w wallet.Wallet, port uint) common.Address {
	tmpClient, err := ethclient.NewEthClient(common.Address{}, "127.0.0.1", port, w, common.Address{})
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
		if err == nil && receipt != nil {
			break
		}
		if !errors.Is(err, ethereum.NotFound) {
			t.Fatal(err)
		}
		t.Logf("Contract deploy tx has not been mined into a block after %s...", time.Since(start))
	}

	t.Logf("Contract deployed to %s - using port %d\n", receipt.ContractAddress, port)
	return receipt.ContractAddress
}
