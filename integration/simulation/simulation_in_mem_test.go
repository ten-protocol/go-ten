package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	ethereum_mock "github.com/obscuronet/go-obscuro/integration/ethereummock"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might give inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupSimTestLog("in-mem")

	numberOfNodes := 7
	numberOfSimWallets := 10
	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)

	simParams := params.SimParams{
		NumberOfNodes:             numberOfNodes,
		AvgBlockDuration:          50 * time.Millisecond,
		SimulationTime:            25 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.5,
		L2ToL1EfficiencyThreshold: 0.5,
		MgmtContractLib:           ethereum_mock.NewMgmtContractLibMock(),
		ERC20ContractLib:          ethereum_mock.NewERC20ContractLibMock(),
		Wallets:                   wallets,
		StartPort:                 integration.StartPortSimulationInMem,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration * 2 / 7

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), &simParams)
}
