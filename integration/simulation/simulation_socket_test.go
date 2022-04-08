package simulation

import (
	"testing"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestSocketNodesMonteCarloSimulation(t *testing.T) {
	logFile := setupTestLog()
	defer logFile.Close()

	params := params.SimParams{
		NumberOfNodes:             10,
		NumberOfWallets:           5,
		AvgBlockDurationUSecs:     uint64(250_000),
		SimulationTimeSecs:        15,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
	}
	params.AvgNetworkLatency = params.AvgBlockDurationUSecs / 15
	params.AvgGossipPeriod = params.AvgBlockDurationUSecs / 3
	params.SimulationTimeUSecs = params.SimulationTimeSecs * 1000 * 1000

	testSimulation(t, network.NewBasicNetworkOfSocketNodes(), params)
}
