package simulation

import (
	"testing"
)

// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes live in the same process, the enclaves run in individual Docker containers, and the Ethereum nodes are mocked out.
func TestDockerNodesMonteCarloSimulation(t *testing.T) {
	params := SimParams{
		NumberOfNodes:         3,
		NumberOfWallets:       5,
		AvgBlockDurationUSecs: uint64(1_000_000),
		SimulationTimeSecs:    15,
	}
	params.AvgNetworkLatency = params.AvgBlockDurationUSecs / 15
	params.AvgGossipPeriod = params.AvgBlockDurationUSecs / 3
	params.SimulationTimeUSecs = params.SimulationTimeSecs * 1000 * 1000

	efficiencies := EfficiencyThresholds{0.2, 0.3, 0.4}

	testSimulation(t, CreateBasicNetworkOfDockerNodes, params, efficiencies)
}
