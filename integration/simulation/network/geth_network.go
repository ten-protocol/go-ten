package network

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"

	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroClients []obscuroclient.Client

	// geth
	gethNetwork  *gethnetwork.GethNetwork
	gethClients  []ethclient.EthClient
	wallets      []wallet.Wallet
	contracts    []string
	workerWallet wallet.Wallet
}

func NewNetworkInMemoryGeth(wallets []wallet.Wallet, workerWallet wallet.Wallet, contracts []string) Network {
	return &networkInMemGeth{
		wallets:      wallets,
		contracts:    contracts,
		workerWallet: workerWallet,
	}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []obscuroclient.Client, []string, error) {
	// kickoff the network with the prefunded wallet addresses
	params.MgmtContractAddr, params.StableTokenContractAddr, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		n.workerWallet,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)

	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.StableTokenContractAddr)

	// Start the obscuro nodes and return the handles
	n.obscuroClients = startInMemoryObscuroNodes(params, stats, n.gethNetwork.GenesisJSON, n.gethClients)

	return n.gethClients, n.obscuroClients, nil, nil
}

func (n *networkInMemGeth) TearDown() {
	// Stop the obscuro nodes first
	StopObscuroNodes(n.obscuroClients)

	// Stop geth last
	StopGethNetwork(n.gethClients, n.gethNetwork)
}
