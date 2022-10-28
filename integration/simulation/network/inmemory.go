package network

import (
	"time"

	"github.com/obscuronet/go-obscuro/integration/datagenerator"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/enclave/bridge"

	"github.com/obscuronet/go-obscuro/integration/simulation/p2p"

	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	"github.com/obscuronet/go-obscuro/integration/ethereummock"
)

type basicNetworkOfInMemoryNodes struct {
	ethNodes       []*ethereummock.Node
	obscuroClients []rpc.Client
}

func NewBasicNetworkOfInMemoryNodes() Network {
	return &basicNetworkOfInMemoryNodes{}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *basicNetworkOfInMemoryNodes) Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	l1Clients := make([]ethadapter.EthClient, params.NumberOfNodes)
	n.ethNodes = make([]*ethereummock.Node, params.NumberOfNodes)
	obscuroNodes := make([]host.MockHost, params.NumberOfNodes)
	n.obscuroClients = make([]rpc.Client, params.NumberOfNodes)
	p2pLayers := make([]*p2p.MockP2P, params.NumberOfNodes)

	// Invent some addresses to assign as the L1 erc20 contracts
	dummyOBXAddress := datagenerator.RandomAddress()
	params.Wallets.Tokens[bridge.HOC].L1ContractAddress = &dummyOBXAddress
	dummyETHAddress := datagenerator.RandomAddress()
	params.Wallets.Tokens[bridge.POC].L1ContractAddress = &dummyETHAddress

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create the in memory l1 and l2 node
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats)
		p2pLayers[i] = p2p.NewMockP2P(params.AvgBlockDuration, params.AvgNetworkLatency)

		agg := createInMemObscuroNode(
			int64(i),
			isGenesis,
			GetNodeType(i),
			params.MgmtContractLib,
			params.ERC20ContractLib,
			params.AvgGossipPeriod,
			stats,
			false,
			nil,
			params.Wallets.NodeWallets[i],
			miner,
			params.Wallets,
			p2pLayers[i],
		)
		obscuroClient := p2p.NewInMemObscuroClient(agg)

		// and connect them to each other
		miner.AddClient(agg)

		n.ethNodes[i] = miner
		obscuroNodes[i] = agg
		n.obscuroClients[i] = obscuroClient
		l1Clients[i] = miner
	}

	// populate the nodes field of each network
	for i := 0; i < params.NumberOfNodes; i++ {
		n.ethNodes[i].Network.(*ethereummock.MockEthNetwork).AllNodes = n.ethNodes
		p2pLayers[i].Nodes = obscuroNodes
	}

	// The sequence of starting the nodes is important to catch various edge cases.
	// Here we first start the mock layer 1 nodes, with a pause between them of a fraction of a block duration.
	// The reason is to make sure that they catch up correctly.
	// Then we pause for a while, to give the L1 network enough time to create a number of blocks, which will have to be ingested by the Obscuro nodes
	// Then, we begin the starting sequence of the Obscuro nodes, again with a delay between them, to test that they are able to cach up correctly.
	// Note: Other simulations might test variations of this pattern.
	for _, m := range n.ethNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 20)
	}

	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	walletClients := createAuthClientsPerWallet(n.obscuroClients, params.Wallets)

	return &RPCHandles{
		EthClients:     l1Clients,
		ObscuroClients: n.obscuroClients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *basicNetworkOfInMemoryNodes) TearDown() {
	for _, client := range n.obscuroClients {
		temp := client
		go func() {
			_ = temp.Call(nil, rpc.StopHost)
			temp.Stop()
		}()
	}

	for _, node := range n.ethNodes {
		temp := node
		go temp.Stop()
	}
}
