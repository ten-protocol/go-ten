package network

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/p2p"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroNodes []*host.Node
	gethNetwork  *gethnetwork.GethNetwork
}

func NewNetworkInMemoryGeth(gethNetwork *gethnetwork.GethNetwork) Network {
	return &networkInMemGeth{
		gethNetwork: gethNetwork,
	}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*host.Node, []string) {
	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	n.obscuroNodes = make([]*host.Node, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		genesis := false
		if i == 0 {
			genesis = true
		}

		// create the in memory l1 and l2 node
		miner := createEthClientConnection(int64(i), n.gethNetwork.WebSocketPorts[i], params.EthWallets[i], params.MgmtContractAddr)
		agg := createInMemObscuroNode(int64(i), genesis, params.TxHandler, params.AvgGossipPeriod, params.AvgBlockDuration, params.AvgNetworkLatency, stats)

		// and connect them to each other
		agg.ConnectToEthNode(miner)

		n.obscuroNodes[i] = agg
		l1Clients[i] = miner
	}

	// make sure the aggregators can talk to each other
	for i := 0; i < params.NumberOfNodes; i++ {
		n.obscuroNodes[i].P2p.(*p2p.MockP2P).Nodes = n.obscuroNodes
	}

	// start each obscuro node
	for _, m := range n.obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroNodes, nil
}

func (n *networkInMemGeth) TearDown() {
	// Nop
}

func createEthClientConnection(id int64, port uint, wallet wallet.Wallet, contractAddr common.Address) ethclient.EthClient {
	ethnode, err := ethclient.NewEthClient(common.BigToAddress(big.NewInt(id)), "127.0.0.1", port, wallet, contractAddr)
	if err != nil {
		panic(err)
	}
	return ethnode
}
