package network

import (
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkOfSocketNodes struct {
	obscuroClients []*obscuroclient.Client
	gethNetwork    *gethnetwork.GethNetwork
}

func NewNetworkOfSocketNodes() Network {
	return &networkOfSocketNodes{}
}

func (n *networkOfSocketNodes) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []*obscuroclient.Client, []string) {
	gethNetwork, contractAddr := createGethNetwork(params)
	n.gethNetwork = &gethNetwork

	params.MgmtContractAddr = contractAddr
	params.TxHandler = mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr)

	l1Clients := make([]ethclient.EthClient, params.NumberOfNodes)
	obscuroNodes := make([]*host.Node, params.NumberOfNodes)
	n.obscuroClients = make([]*obscuroclient.Client, params.NumberOfNodes)
	nodeP2pAddrs := make([]string, params.NumberOfNodes)

	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign a P2P address to each node on the network.
		nodeP2pAddrs[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+200+i)
	}

	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		// create a remote enclave server
		nodeID := common.BigToAddress(big.NewInt(int64(i)))
		enclaveAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+300+i)
		_, err := enclave.StartServer(enclaveAddr, nodeID, params.TxHandler, false, nil, stats)
		if err != nil {
			panic(fmt.Sprintf("failed to create enclave server: %v", err))
		}

		// create the L1 client, the Obscuro host and enclave service, and the L2 client
		l1Client := createEthClientConnection(
			int64(i),
			n.gethNetwork.WebSocketPorts[i],
			params.EthWallets[i],
			params.MgmtContractAddr,
		)
		obscuroClientAddr := fmt.Sprintf("%s:%d", Localhost, params.StartPort+400+i)
		obscuroClient := obscuroclient.NewClient(obscuroClientAddr)
		agg := createSocketObscuroNode(int64(i), isGenesis, params.AvgGossipPeriod, stats, nodeP2pAddrs[i], nodeP2pAddrs, enclaveAddr, obscuroClientAddr, params.TxHandler)

		// connect the L1 and L2 nodes
		agg.ConnectToEthNode(l1Client)

		obscuroNodes[i] = agg
		n.obscuroClients[i] = &obscuroClient
		l1Clients[i] = l1Client
	}

	// start each obscuro node
	for _, m := range obscuroNodes {
		t := m
		go t.Start()
		time.Sleep(params.AvgBlockDuration / 3)
	}

	return l1Clients, n.obscuroClients, nodeP2pAddrs
}

func (n *networkOfSocketNodes) TearDown() {
	defer n.gethNetwork.StopNodes()

	for _, client := range n.obscuroClients {
		temp := client
		go func() {
			defer (*temp).Stop()
			(*temp).Call(nil, obscuroclient.RPCStopHost) //nolint:errcheck
		}()
	}
}
