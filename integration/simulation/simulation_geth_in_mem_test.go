package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// TestInMemoryGethMonteCarloSimulation runs the simulation against a private geth network using Clique (PoA)
func TestInMemoryGethMonteCarloSimulation(t *testing.T) {
	setupTestLog("geth-in-mem")

	numberOfNodes := 5

	// there is one wallet per node, so there have to be at least numberOfNodes wallets available
	numberOfWallets := numberOfNodes

	// randomly create the ethereum wallets to be used and prefund them
	wallets := make([]wallet.Wallet, numberOfWallets)
	for i := 0; i < numberOfWallets; i++ {
		wallets[i] = datagenerator.RandomWallet()
	}

	simParams := params.SimParams{
		NumberOfNodes:             numberOfNodes,
		NumberOfObscuroWallets:    numberOfWallets,
		AvgBlockDuration:          6 * time.Second,
		SimulationTime:            60 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.5,
		L2ToL1EfficiencyThreshold: 0.5, // one rollup every 2 blocks
		EthWallets:                wallets,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, network.NewNetworkInMemoryGeth(), &simParams)
}
