package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process. The L1 network is a private geth network using Clique (PoA).
func TestFullNetworkMonteCarloSimulation(t *testing.T) {
	setupTestLog("full-network")

	numberOfNodes := 5
	numberOfSimWallets := 5

	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, 1, integration.EthereumChainID, integration.ObscuroChainID)

	simParams := &params.SimParams{
		NumberOfNodes:         numberOfNodes,
		AvgBlockDuration:      1 * time.Second,
		SimulationTime:        35 * time.Second,
		L1EfficiencyThreshold: 0.2,
		// Very hard to have precision here as blocks are continually produced and not dependent on the simulation execution thread
		L2EfficiencyThreshold:     0.6, // nodes might produce rollups because they receive a new block
		L2ToL1EfficiencyThreshold: 0.7, // nodes might stop producing rollups but the geth network is still going
		Wallets:                   wallets,
		StartPort:                 integration.StartPortSimulationFullNetwork,
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, network.NewNetworkOfSocketNodes(wallets), simParams)
}
