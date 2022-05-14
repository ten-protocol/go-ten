package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/contracts"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// TestGethMemObscuroEthERC20MonteCarloSimulation runs the simulation against a private geth network using Clique (PoA)
func TestGethMemObscuroEthERC20MonteCarloSimulation(t *testing.T) {
	setupTestLog()

	numberOfNodes := 5

	// there is one wallet per node, so there have to be at least numberOfNodes wallets available
	numberOfWallets := numberOfNodes

	// randomly create the ethereum wallets to be used and prefund them
	// create one extra wallet as the worker wallet
	wallets := make([]wallet.Wallet, numberOfWallets+1)
	for i := 0; i < numberOfWallets+1; i++ {
		wallets[i] = datagenerator.RandomWallet()
	}

	// The last wallet as the worker wallet ( to deposit and inject transactions )
	workerWallet := wallets[numberOfWallets]

	// define contracts to be deployed
	contractsBytes := []string{
		contracts.MgmtContractByteCode,
		contracts.PedroERC20ContractByteCode,
	}

	// define the network to use
	netw := network.NewNetworkInMemoryGeth(wallets, workerWallet, contractsBytes)

	simParams := &params.SimParams{
		NumberOfNodes:             numberOfNodes,
		NumberOfObscuroWallets:    numberOfWallets,
		AvgBlockDuration:          1 * time.Second,
		SimulationTime:            60 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.9,
		L2ToL1EfficiencyThreshold: 0.9, // one rollup every 2 blocks
		EthWallets:                wallets,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, netw, simParams)
}
