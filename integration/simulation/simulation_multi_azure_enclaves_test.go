////go:build azure
//// +build azure
//
package simulation

//
//import (
//	"testing"
//	"time"
//
//	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
//
//	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
//)
//
//// TODO: we really need tests to demonstrate the unhappy-cases in the attestation scenario:
////		- if someone puts a dodgy public key on a request with a genuine attestation report they shouldn't get secret
////		- if owner doesn't match - they shouldn't get secret
//
//// Todo: replace with the IPs of the VMs you are testing, see the azuredeployer README for more info.
////		If you are adding more than two IPs, be sure to increase the NumberOfNodes in the SimParams below.
//var vmIPs = []string{"20.254.65.172", "20.254.67.124"}
//
//// This test creates a network of L2 nodes consisting of just the Azure nodes configured above.
////
//// It then injects transactions, and finally checks the resulting output blockchain
//// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
//// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
//func TestMultipleAzureEnclaveNodesVerifyEachOther(t *testing.T) {
//	setupTestLog()
//
//	params := params.SimParams{
//		NumberOfNodes:             2,
//		NumberOfObscuroWallets:    5,
//		AvgBlockDuration:          time.Second,
//		SimulationTime:            30 * time.Second,
//		L1EfficiencyThreshold:     0.2,
//		L2EfficiencyThreshold:     0.3,
//		L2ToL1EfficiencyThreshold: 0.4,
//	}
//	params.AvgNetworkLatency = params.AvgBlockDuration / 15
//	params.AvgGossipPeriod = params.AvgBlockDuration / 3
//
//	enclaves := make([]string, len(vmIPs))
//	for i, ip := range vmIPs {
//		enclaves[i] = ip + ":11000"
//	}
//	testSimulation(t, network.NewNetworkWithAzureEnclaves(enclaves), &params)
//}
