package simulation

//
//import (
//	"testing"
//	"time"
//
//	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
//	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
//
//	"github.com/obscuronet/obscuro-playground/integration"
//
//	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
//
//	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
//)
//
//// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
//// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
//// All nodes and enclaves live in the same process. The L1 network is a private geth network using Clique (PoA).
//func TestFullNetworkMonteCarloSimulation(t *testing.T) {
//	setupTestLog("socket")
//
//	numberOfNodes := 5
//	numberOfWallets := numberOfNodes // We need at least one wallet per node.
//
//	// randomly create the ethereum wallets to be used and prefund them
//	wallets := make([]wallet.Wallet, numberOfWallets)
//	for i := 0; i < numberOfWallets; i++ {
//		wallets[i] = datagenerator.RandomWallet()
//	}
//
//	simParams := params.SimParams{
//		NumberOfNodes:             numberOfNodes,
//		NumberOfObscuroWallets:    numberOfWallets,
//		AvgBlockDuration:          6 * time.Second,
//		SimulationTime:            60 * time.Second,
//		L1EfficiencyThreshold:     0.2,
//		L2EfficiencyThreshold:     0.5,
//		L2ToL1EfficiencyThreshold: 0.5, // one rollup every 2 blocks
//		EthWallets:                wallets,
//		StartPort:                 integration.StartPortSimulationSocket,
//	}
//	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
//	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3
//
//	testSimulation(t, network.NewNetworkOfSocketNodes(), &simParams)
//}
