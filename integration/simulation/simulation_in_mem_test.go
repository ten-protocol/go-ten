package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/ethereummock"

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

	simParams := params.SimParams{
		NumberOfNodes:             10,
		NumberOfWallets:           5,
		AvgBlockDuration:          50 * time.Millisecond,
		SimulationTime:            15 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.32,
		L2ToL1EfficiencyThreshold: 0.34,
		TxHandler:                 ethereummock.NewMockTxHandler(),
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration * 2 / 7

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), simParams)
}
