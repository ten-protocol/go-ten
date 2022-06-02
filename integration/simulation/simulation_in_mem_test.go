package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	ethereum_mock "github.com/obscuronet/go-obscuro/integration/ethereummock"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might result in inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupTestLog("in-mem")

	simParams := params.SimParams{
		NumberOfNodes:             7,
		AvgBlockDuration:          50 * time.Millisecond,
		SimulationTime:            25 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.5,
		L2ToL1EfficiencyThreshold: 0.5,
		MgmtContractLib:           ethereum_mock.NewMgmtContractLibMock(),
		ERC20ContractLib:          ethereum_mock.NewERC20ContractLibMock(),
		StartPort:                 integration.StartPortSimulationInMem,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration * 2 / 7

	for i := 0; i < simParams.NumberOfNodes+1; i++ {
		simParams.NodeEthWallets = append(simParams.NodeEthWallets, datagenerator.RandomWallet(integration.EthereumChainID))
		simParams.SimEthWallets = append(simParams.SimEthWallets, datagenerator.RandomWallet(integration.EthereumChainID))
	}

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), &simParams)
}
