package network

import (
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

type networkInMemGeth struct {
	l2Clients []rpc.Client

	// geth
	gethNetwork *gethnetwork.GethNetwork
	gethClients []ethadapter.EthClient
	wallets     *params.SimWallets
}

func NewNetworkInMemoryGeth(wallets *params.SimWallets) Network {
	return &networkInMemGeth{
		wallets: wallets,
	}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	// kickoff the network with the prefunded wallet addresses
	params.MgmtContractAddr, params.ObxErc20Address, params.EthErc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)

	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr, testlog.Logger())
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.ObxErc20Address, params.EthErc20Address)

	// Start the obscuro nodes and return the handles
	n.l2Clients = startInMemoryObscuroNodes(params, stats, n.gethNetwork.GenesisJSON, n.gethClients)

	obscuroClients := make([]*obsclient.ObsClient, params.NumberOfNodes)
	for idx, l2Client := range n.l2Clients {
		obscuroClients[idx] = obsclient.NewObsClient(l2Client)
	}
	walletClients := createAuthClientsPerWallet(n.l2Clients, params.Wallets)

	return &RPCHandles{
		EthClients:     n.gethClients,
		ObscuroClients: obscuroClients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *networkInMemGeth) TearDown() {
	// Stop the obscuro nodes first
	StopObscuroNodes(n.l2Clients)

	// Stop geth last
	StopGethNetwork(n.gethClients, n.gethNetwork)
}
