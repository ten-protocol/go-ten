package network

import (
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkWithAzureEnclaves struct {
	ethNodes         []*ethereum_mock.Node
	obscuroNodes     []*host.Node
	obscuroClients   []*obscuroclient.Client
	enclaveAddresses []string
}

func NewNetworkWithOneAzureEnclave(enclaveAddress string) Network {
	return &networkWithAzureEnclaves{enclaveAddresses: []string{enclaveAddress}}
}

func NewNetworkWithAzureEnclaves(enclaveAddresses []string) Network {
	return &networkWithAzureEnclaves{enclaveAddresses: enclaveAddresses}
}

func (n *networkWithAzureEnclaves) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string, error) {
	if len(n.enclaveAddresses) == 0 {
		panic("Cannot create azure enclaves network without at least one enclave address.")
	}

	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	n.ethNodes = make([]*ethereum_mock.Node, params.NumberOfNodes)
	n.obscuroNodes = make([]*host.Node, params.NumberOfNodes)
	n.obscuroClients = make([]*obscuroclient.Client, params.NumberOfNodes)
	nodeP2pAddrs := make([]string, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeP2pAddrs[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+i)
	}

	// set up nodes with azure enclave
	for i := 0; i < len(n.enclaveAddresses); i++ {
		isGenesis := i == 0
		// create the in memory l1 and l2 node
		obscuroClientAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+200+i)
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats)
		agg := createSocketObscuroNode(
			int64(i),
			isGenesis,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			n.enclaveAddresses[i],
			obscuroClientAddr,
			params.NodeEthWallets[i],
			params.MgmtContractLib,
		)
		obscuroClient := obscuroclient.NewClient(obscuroClientAddr)

		n.wireUpNode(i, l1Clients, miner, agg, &obscuroClient)
	}

	// set up nodes with mock enclaves
	for i := len(n.enclaveAddresses); i < params.NumberOfNodes; i++ {
		// create a remote enclave server
		enclavePort := uint64(params.StartPort + DefaultWsPortOffset + i)
		enclaveAddress := fmt.Sprintf("%s:%d", Localhost, enclavePort)
		enclaveConfig := config.EnclaveConfig{
			HostID:           common.BigToAddress(big.NewInt(int64(i))),
			Address:          fmt.Sprintf("%s:%d", Localhost, enclavePort),
			L1ChainID:        integration.EthereumChainID,
			ObscuroChainID:   integration.ObscuroChainID,
			WillAttest:       false,
			ValidateL1Blocks: false,
			GenesisJSON:      nil,
			UseInMemoryDB:    true,
		}
		_, err := enclave.StartServer(
			enclaveConfig,
			params.MgmtContractLib,
			params.ERC20ContractLib,
			stats,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}

		// create the in memory l1 and l2 node
		obscuroClientAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+200+i)
		miner := createMockEthNode(int64(i), params.NumberOfNodes, params.AvgBlockDuration, params.AvgNetworkLatency, stats)
		agg := createSocketObscuroNode(
			int64(i),
			false,
			params.AvgGossipPeriod,
			stats,
			nodeP2pAddrs[i],
			nodeP2pAddrs,
			enclaveAddress,
			obscuroClientAddr,
			params.NodeEthWallets[i],
			params.MgmtContractLib,
		)
		obscuroClient := obscuroclient.NewClient(obscuroClientAddr)

		n.wireUpNode(i, l1Clients, miner, agg, &obscuroClient)
	}

	// populate the nodes field of the L1 network
	for i := 0; i < params.NumberOfNodes; i++ {
		n.ethNodes[i].Network.(*ethereum_mock.MockEthNetwork).AllNodes = n.ethNodes
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
		time.Sleep(params.AvgBlockDuration / 8)
	}

	time.Sleep(params.AvgBlockDuration * 20)
	for _, m := range n.obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroClients, nodeP2pAddrs, nil
}

func (n *networkWithAzureEnclaves) TearDown() {
	for _, client := range n.obscuroClients {
		temp := client
		go (*temp).Stop()
	}

	for _, node := range n.ethNodes {
		temp := node
		go temp.Stop()
	}
}

func (n *networkWithAzureEnclaves) wireUpNode(idx int, l1Clients []ethclient.EthClient, miner *ethereum_mock.Node, agg *host.Node, obscuroClient *obscuroclient.Client) {
	// and connect them to each other
	agg.ConnectToEthNode(miner)
	miner.AddClient(agg)

	n.ethNodes[idx] = miner
	n.obscuroNodes[idx] = agg
	n.obscuroClients[idx] = obscuroClient
	l1Clients[idx] = miner
}
