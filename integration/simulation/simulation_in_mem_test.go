package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration"
	ethereum_mock "github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might give inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupSimTestLog("in-mem")

	// todo (#718) - try increasing this back to 7 once faster-finality model is optimised
	numberOfNodes := 5
	numberOfSimWallets := 10
	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)

	simParams := params.SimParams{
		NumberOfNodes: numberOfNodes,
		//  todo (#718) - try reducing this back to 50 milliseconds once faster-finality model is optimised
		AvgBlockDuration:          100 * time.Millisecond,
		SimulationTime:            30 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.5,
		L2ToL1EfficiencyThreshold: 0.5,
		MgmtContractLib:           ethereum_mock.NewMgmtContractLibMock(),
		ERC20ContractLib:          ethereum_mock.NewERC20ContractLibMock(),
		Wallets:                   wallets,
		StartPort:                 integration.StartPortSimulationInMem,
		IsInMem:                   true,
		L1SetupData:               &params.L1SetupData{},
		ReceiptTimeout:            5 * time.Second,
		StoppingDelay:             2 * time.Second,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), &simParams)
}
