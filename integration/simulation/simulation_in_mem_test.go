package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/buildhelper/helpertypes"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might result in inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupTestLog()

	params := params.SimParams{
		NumberOfNodes:             10,
		NumberOfWallets:           5,
		AvgBlockDuration:          40 * time.Microsecond,
		SimulationTime:            15 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.32,
		L2ToL1EfficiencyThreshold: 0.5,
	}

	params.AvgNetworkLatency = params.AvgBlockDuration / 15
	params.AvgGossipPeriod = params.AvgBlockDuration / 3

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), params)
}

func TestMemObscuroRealEthMonteCarloSimulation(t *testing.T) {
	t.Skip("test under construction")
	setupTestLog()

	helpertypes.IsRealEth = true
	params := params.SimParams{
		NumberOfNodes:             2,
		NumberOfWallets:           2,
		AvgBlockDuration:          40 * time.Microsecond,
		SimulationTime:            15 * time.Second,
		L1EfficiencyThreshold:     0.9, // todo review this
		L2EfficiencyThreshold:     0.9,
		L2ToL1EfficiencyThreshold: 0.9,
	}

	params.AvgNetworkLatency = params.AvgBlockDuration / 15
	params.AvgGossipPeriod = params.AvgBlockDuration / 2

	testSimulation(t, network.NewNetworkInMemoryGeth(), params)
}
