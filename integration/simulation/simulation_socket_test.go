package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestSocketNodesMonteCarloSimulation(t *testing.T) {
	setupTestLog()

	params := params.SimParams{
		NumberOfNodes:             10,
		NumberOfWallets:           5,
		AvgBlockDuration:          250 * time.Microsecond,
		SimulationTime:            15 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
	}
	params.AvgNetworkLatency = params.AvgBlockDuration / 15
	params.AvgGossipPeriod = params.AvgBlockDuration / 3

	testSimulation(t, network.NewBasicNetworkOfSocketNodes(), params)
}
