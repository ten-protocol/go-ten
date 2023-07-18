package network

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/host/container"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	testcommon "github.com/obscuronet/go-obscuro/integration/common"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/p2p"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

type basicNetworkOfInMemoryNodes struct {
	ethNodes  []*ethereummock.Node
	l2Clients []rpc.Client
}

func NewBasicNetworkOfInMemoryNodes() Network {
	return &basicNetworkOfInMemoryNodes{}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *basicNetworkOfInMemoryNodes) Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	l1Clients := make([]ethadapter.EthClient, params.NumberOfNodes)
	n.ethNodes = make([]*ethereummock.Node, params.NumberOfNodes)
	obscuroNodes := make([]*container.HostContainer, params.NumberOfNodes)
	n.l2Clients = make([]rpc.Client, params.NumberOfNodes)
	obscuroHosts := make([]host.Host, params.NumberOfNodes)

	p2pNetw := p2p.NewMockP2PNetwork(params.AvgBlockDuration, params.AvgNetworkLatency)

	// Invent some addresses to assign as the L1 erc20 contracts
	dummyOBXAddress := datagenerator.RandomAddress()
	params.Wallets.Tokens[testcommon.HOC].L1ContractAddress = &dummyOBXAddress
	dummyETHAddress := datagenerator.RandomAddress()
	params.Wallets.Tokens[testcommon.POC].L1ContractAddress = &dummyETHAddress
	disabledBus := common.BigToAddress(common.Big0)

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create the in memory l1 and l2 node
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats)

		agg := createInMemObscuroNode(
			int64(i),
			isGenesis,
			GetNodeType(i),
			params.MgmtContractLib,
			false,
			nil,
			params.Wallets.NodeWallets[i],
			miner,
			p2pNetw.NewNode(i),
			&disabledBus,
			common.Hash{},
			params.AvgBlockDuration/2,
		)
		obscuroClient := p2p.NewInMemObscuroClient(agg)

		n.ethNodes[i] = miner
		obscuroNodes[i] = agg
		n.l2Clients[i] = obscuroClient
		l1Clients[i] = miner
		obscuroHosts[i] = obscuroNodes[i].Host()
	}

	// populate the nodes field of each network
	for i := 0; i < params.NumberOfNodes; i++ {
		n.ethNodes[i].Network.(*ethereummock.MockEthNetwork).AllNodes = n.ethNodes
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
		time.Sleep(params.AvgBlockDuration)
	}

	for _, m := range obscuroNodes {
		t := m
		go func() {
			err := t.Start()
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	obscuroClients := make([]*obsclient.ObsClient, params.NumberOfNodes)
	for idx, l2Client := range n.l2Clients {
		obscuroClients[idx] = obsclient.NewObsClient(l2Client)
	}
	walletClients := createAuthClientsPerWallet(n.l2Clients, params.Wallets)

	return &RPCHandles{
		EthClients:     l1Clients,
		ObscuroClients: obscuroClients,
		RPCClients:     n.l2Clients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *basicNetworkOfInMemoryNodes) TearDown() {
	StopObscuroNodes(n.l2Clients)

	for _, node := range n.ethNodes {
		temp := node
		go temp.Stop()
	}
}
