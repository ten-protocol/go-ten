package simulation

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

const (
	localhost        = "localhost"
	p2pStartPort     = 10000
	enclaveStartPort = 11000
	testLogs         = "../.build/simulations/"
)

func setupTestLog() *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(testLogs, 0o700)
	if err != nil {
		panic(err)
	}
	f, err := os.CreateTemp(testLogs, fmt.Sprintf("simulation-result-%d-*.txt", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	log.SetLog(f)
	return f
}

func createMockEthNode(id int64, nrNodes int, avgBlockDurationUSecs uint64, avgNetworkLatency uint64, stats *Stats) *ethereum_mock.Node {
	mockEthNetwork := ethereum_mock.NewMockEthNetwork(avgBlockDurationUSecs, avgNetworkLatency, stats)
	ethereumMockCfg := defaultMockEthNodeCfg(nrNodes, avgBlockDurationUSecs)
	// create an in memory mock ethereum node responsible with notifying the layer 2 node about blocks
	miner := ethereum_mock.NewMiner(common.BigToAddress(big.NewInt(id)), ethereumMockCfg, mockEthNetwork, stats)
	mockEthNetwork.CurrentNode = &miner
	return &miner
}

func createInMemObscuroNode(id int64, genesis bool, avgGossipPeriod uint64, avgBlockDurationUSecs uint64, avgNetworkLatency uint64, stats *Stats) *host.Node {
	obscuroInMemNetwork := NewMockP2P(avgBlockDurationUSecs, avgNetworkLatency)

	obscuroNodeCfg := defaultObscuroNodeCfg(avgGossipPeriod)

	nodeID := common.BigToAddress(big.NewInt(id))
	enclaveClient := enclave.NewEnclave(nodeID, true, stats)

	// create an in memory obscuro node
	node := host.NewObscuroAggregator(nodeID, obscuroNodeCfg, nil, stats, genesis, enclaveClient, obscuroInMemNetwork)
	obscuroInMemNetwork.currentNode = &node
	return &node
}

func createSocketObscuroNode(id int64, genesis bool, avgGossipPeriod uint64, stats *Stats, peerAddrs []string, enclavePort uint64) *host.Node {
	nodeID := common.BigToAddress(big.NewInt(id))

	// create an enclave client
	enclaveAddr := fmt.Sprintf("%s:%d", localhost, enclavePort)
	enclaveClient := host.NewEnclaveRPCClient(enclaveAddr, host.ClientRPCTimeoutSecs*time.Second, nodeID)

	// create a socket obscuro node
	nodeP2p := p2p.NewSocketP2PLayer(peerAddrs[id-1], peerAddrs)
	obscuroNodeCfg := defaultObscuroNodeCfg(avgGossipPeriod)
	node := host.NewObscuroAggregator(nodeID, obscuroNodeCfg, nil, stats, genesis, enclaveClient, nodeP2p)

	return &node
}

// creates the nodes, wires them up, and populates the network objects
func CreateBasicNetworkOfInMemoryNodes(params SimParams, stats *Stats) ([]*ethereum_mock.Node, []*host.Node) {
	// todo - add observer nodes
	l1Nodes := make([]*ethereum_mock.Node, params.NumberOfNodes)
	l2Nodes := make([]*host.Node, params.NumberOfNodes)
	for i := 1; i <= params.NumberOfNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}

		// create the in memory l1 and l2 node
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDurationUSecs, params.AvgNetworkLatency, stats)
		agg := createInMemObscuroNode(int64(i), genesis, params.AvgGossipPeriod, params.AvgBlockDurationUSecs, params.AvgNetworkLatency, stats)

		// and connect them to each other
		agg.ConnectToEthNode(miner)
		miner.AddClient(agg)

		l1Nodes[i-1] = miner
		l2Nodes[i-1] = agg
	}

	// populate the nodes field of each network
	for i := 0; i < params.NumberOfNodes; i++ {
		l1Nodes[i].Network.(*ethereum_mock.MockEthNetwork).AllNodes = l1Nodes
		l2Nodes[i].P2p.(*MockP2P).Nodes = l2Nodes
	}

	return l1Nodes, l2Nodes
}

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
func CreateBasicNetworkOfSocketNodes(params SimParams, stats *Stats) ([]*ethereum_mock.Node, []*host.Node) {
	// todo - add observer nodes
	l1Nodes := make([]*ethereum_mock.Node, params.NumberOfNodes)
	l2Nodes := make([]*host.Node, params.NumberOfNodes)

	var nodeAddrs []string
	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeAddrs = append(nodeAddrs, fmt.Sprintf("%s:%d", localhost, p2pStartPort+i))
	}

	for i := 1; i <= params.NumberOfNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}

		// create a remote enclave server
		nodeID := common.BigToAddress(big.NewInt(int64(i)))
		enclavePort := uint64(enclaveStartPort + i)
		enclaveAddress := fmt.Sprintf("localhost:%d", enclavePort)
		err := enclave.StartServer(enclaveAddress, nodeID, stats)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}

		// create the in memory l1 and l2 node
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDurationUSecs, params.AvgNetworkLatency, stats)
		agg := createSocketObscuroNode(int64(i), genesis, params.AvgGossipPeriod, stats, nodeAddrs, enclavePort)

		// and connect them to each other
		agg.ConnectToEthNode(miner)
		miner.AddClient(agg)

		l1Nodes[i-1] = miner
		l2Nodes[i-1] = agg
	}

	// populate the nodes field of the L1 network
	for i := 0; i < params.NumberOfNodes; i++ {
		l1Nodes[i].Network.(*ethereum_mock.MockEthNetwork).AllNodes = l1Nodes
	}

	return l1Nodes, l2Nodes
}

// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.
// creates Obscuro nodes with their own Dockerised enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
func CreateBasicNetworkOfDockerNodes(params SimParams, stats *Stats) ([]*ethereum_mock.Node, []*host.Node) {
	// todo - add observer nodes
	l1Nodes := make([]*ethereum_mock.Node, params.NumberOfNodes)
	l2Nodes := make([]*host.Node, params.NumberOfNodes)

	var nodeAddrs []string
	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeAddrs = append(nodeAddrs, fmt.Sprintf("%s:%d", localhost, p2pStartPort+i))
	}

	for i := 1; i <= params.NumberOfNodes; i++ {
		genesis := false
		if i == 1 {
			genesis = true
		}

		// create the in memory l1 and l2 node
		enclavePort := uint64(enclaveStartPort + i - 1)
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDurationUSecs, params.AvgNetworkLatency, stats)
		agg := createSocketObscuroNode(int64(i), genesis, params.AvgGossipPeriod, stats, nodeAddrs, enclavePort)

		// and connect them to each other
		agg.ConnectToEthNode(miner)
		miner.AddClient(agg)

		l1Nodes[i-1] = miner
		l2Nodes[i-1] = agg
	}

	// populate the nodes field of the L1 network
	for i := 0; i < params.NumberOfNodes; i++ {
		l1Nodes[i].Network.(*ethereum_mock.MockEthNetwork).AllNodes = l1Nodes
	}

	return l1Nodes, l2Nodes
}

func defaultObscuroNodeCfg(gossipPeriod uint64) host.AggregatorCfg {
	return host.AggregatorCfg{ClientRPCTimeoutSecs: host.ClientRPCTimeoutSecs, GossipRoundDuration: gossipPeriod}
}

func defaultMockEthNodeCfg(nrNodes int, avgBlockDuration uint64) ethereum_mock.MiningConfig {
	return ethereum_mock.MiningConfig{
		PowTime: func() uint64 {
			// This formula might feel counter-intuitive, but it is a good approximation for Proof of Work.
			// It creates a uniform distribution up to nrMiners*avgDuration
			// Which means on average, every round, the winner (miner who gets the lowest nonce) will pick a number around "avgDuration"
			// while everyone else will have higher values.
			// Over a large number of rounds, the actual average block duration will be around the desired value, while the number of miners who get very close numbers will be limited.
			return obscurocommon.RndBtw(avgBlockDuration/uint64(nrNodes), uint64(nrNodes)*avgBlockDuration)
		},
	}
}

func minMax(arr []uint64) (min uint64, max uint64) {
	min = ^uint64(0)
	for _, no := range arr {
		if no < min {
			min = no
		}
		if no > max {
			max = no
		}
	}
	return
}
