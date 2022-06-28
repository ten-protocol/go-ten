package network

import (
	"github.com/obscuronet/obscuro-playground/go/ethadapter"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroClients []rpcclientlib.Client

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
func (n *networkInMemGeth) Create(params *params.SimParams, stats *stats.Stats) ([]ethadapter.EthClient, []rpcclientlib.Client, error) {
	// kickoff the network with the prefunded wallet addresses
	params.MgmtContractAddr, params.BtcErc20Address, params.EthErc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)

	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.BtcErc20Address, params.EthErc20Address)

	// Start the obscuro nodes and return the handles
	n.obscuroClients = startInMemoryObscuroNodes(params, stats, n.gethNetwork.GenesisJSON, n.gethClients)

	return n.gethClients, n.obscuroClients, nil
}

func (n *networkInMemGeth) TearDown() {
	// Stop the obscuro nodes first
	StopObscuroNodes(n.obscuroClients)

	// Stop geth last
	StopGethNetwork(n.gethClients, n.gethNetwork)
}
