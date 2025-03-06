package network

import (
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common/retry"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/host/container"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	testcommon "github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/p2p"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/integration/simulation/stats"
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
	tenNodes := make([]*container.HostContainer, params.NumberOfNodes)
	n.l2Clients = make([]rpc.Client, params.NumberOfNodes)
	tenHosts := make([]host.Host, params.NumberOfNodes)

	p2pNetw := p2p.NewMockP2PNetwork(params.AvgBlockDuration, params.AvgNetworkLatency, params.NodeWithInboundP2PDisabled)

	// Invent some addresses to assign as the L1 erc20 contracts
	dummyOBXAddress := datagenerator.RandomAddress()
	params.Wallets.Tokens[testcommon.HOC].L1ContractAddress = &dummyOBXAddress
	dummyETHAddress := datagenerator.RandomAddress()
	params.Wallets.Tokens[testcommon.POC].L1ContractAddress = &dummyETHAddress
	dummyBus := datagenerator.RandomAddress()
	// dummyMgmtContractAddress := datagenerator.RandomAddress()
	// params.MgmtContractLib

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		incomingP2PDisabled := !isGenesis && i == params.NodeWithInboundP2PDisabled

		// create the in memory l1 and l2 node
		miner := createMockEthNode(i, params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats, params.BlobResolver)
		agg := createInMemTenNode(
			int64(i),
			isGenesis,
			GetNodeType(i),
			params.NetworkContractConfigLib,
			params.Wallets.NodeWallets[i],
			miner,
			p2pNetw.NewNode(i),
			dummyBus,
			common.Hash{},
			params.AvgBlockDuration/2,
			incomingP2PDisabled,
			params.AvgBlockDuration,
			params.BlobResolver,
		)
		tenClient := p2p.NewInMemTenClient(agg)

		n.ethNodes[i] = miner
		tenNodes[i] = agg
		n.l2Clients[i] = tenClient
		l1Clients[i] = miner
		tenHosts[i] = tenNodes[i].Host()
	}

	// populate the nodes field of each network
	for i := 0; i < params.NumberOfNodes; i++ {
		n.ethNodes[i].Network.(*ethereummock.MockEthNetwork).AllNodes = n.ethNodes
	}

	// The sequence of starting the nodes is important to catch various edge cases.
	// Here we first start the mock layer 1 nodes, with a pause between them of a fraction of a block duration.
	// The reason is to make sure that they catch up correctly.
	// Then we pause for a while, to give the L1 network enough time to create a number of blocks, which will have to be ingested by the en nodes
	// Then, we begin the starting sequence of the TEN nodes, again with a delay between them, to test that they are able to cach up correctly.
	// Note: Other simulations might test variations of this pattern.
	for _, m := range n.ethNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration)
	}

	for _, m := range tenNodes {
		t := m
		go func() {
			err := t.Start()
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	tenClients := make([]*obsclient.ObsClient, params.NumberOfNodes)
	for idx, l2Client := range n.l2Clients {
		tenClients[idx] = obsclient.NewObsClient(l2Client)
	}
	walletClients := createAuthClientsPerWallet(n.l2Clients, params.Wallets)

	var sequencerHealth host.HealthCheck
	// wait for the sequencer to be healthy
	err := retry.Do(func() error {
		var err error
		sequencerHealth, err = tenClients[0].Health()
		if err != nil {
			return err
		}
		if len(sequencerHealth.Enclaves) == 0 {
			return fmt.Errorf("no enclaves available to promote on sequencer")
		}

		// the nodes are healthy, we can continue
		return nil
	}, retry.NewTimeoutStrategy(30*params.AvgBlockDuration, params.AvgBlockDuration))
	if err != nil {
		panic(err)
	}

	// mock implementation of the permissioning, tell the mock L1 the seq address
	for _, node := range n.ethNodes {
		node.PromoteEnclave(sequencerHealth.Enclaves[0].EnclaveID)
	}
	permMockAddr := ethereummock.GrantSeqTxAddr
	mockTx := types.NewTx(&types.LegacyTx{
		To:   &permMockAddr,
		Data: []byte{0x1},
	})
	err = n.ethNodes[0].SendTransaction(mockTx)
	if err != nil {
		return nil, err
	}

	return &RPCHandles{
		EthClients:     l1Clients,
		TenClients:     tenClients,
		RPCClients:     n.l2Clients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *basicNetworkOfInMemoryNodes) TearDown() {
	StopTenNodes(n.l2Clients)

	for _, node := range n.ethNodes {
		temp := node
		go temp.Stop()
	}
}
