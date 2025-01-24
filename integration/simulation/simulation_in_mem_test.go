package simulation

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might give inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupSimTestLog("in-mem")

	numberOfNodes := 5
	numberOfSimWallets := 10
	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.TenChainID)

	simParams := params.SimParams{
		NumberOfNodes:              numberOfNodes,
		AvgBlockDuration:           180 * time.Millisecond,
		SimulationTime:             45 * time.Second,
		L1EfficiencyThreshold:      0.8,
		MgmtContractLib:            ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib:           ethereummock.NewERC20ContractLibMock(),
		BlobResolver:               ethereummock.NewMockBlobResolver(),
		Wallets:                    wallets,
		StartPort:                  integration.TestPorts.TestInMemoryMonteCarloSimulationPort,
		IsInMem:                    true,
		L1TenData:                  &params.L1TenData{},
		ReceiptTimeout:             6 * time.Second,
		StoppingDelay:              8 * time.Second,
		NodeWithInboundP2PDisabled: 2,
		L1BeaconPort:               integration.TestPorts.TestInMemoryMonteCarloSimulationPort + integration.DefaultPrysmGatewayPortOffset,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), &simParams)
}
