package network

import (
	"math/rand"

	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

// Network is responsible with knowing how to manage the lifecycle of networks of Ethereum or Obscuro nodes.
// These networks can be composed of in-memory go-routines or of fully fledged existing nodes like Ropsten.
// Implementation notes:
// - This is a work in progress, so there is a lot of code duplication in the implementations
// - Once we implement a few more versions: for example using Geth, we'll revisit and create better abstractions.
// TODO Decompose the network so we can pick and choose different types of l1 and obscuro nodes
type Network interface {
	// Create - returns the started Ethereum nodes and the started Obscuro node clients.
	// Responsible with spinning up all resources required for the test
	// Return an error in case it cannot start for an expected reason. Otherwise it panics.
	Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error)
	TearDown()
}

type RPCHandles struct {
	// an eth client per eth node in the network
	EthClients []ethadapter.EthClient

	// an obscuro client per obscuro node in the network (used for things like validation rather than transactions on behalf of sim accounts)
	ObscuroClients []rpc.Client

	// an RPC client per node per wallet, with a viewing key set up (on the client and registered on its corresponding host enclave),
	//	to mimic user acc interaction via a wallet extension
	// map of owner addresses to RPC clients for that owner (one per L2 node)
	// todo: simplify this with a client per node when we have clients that can support multiple wallets
	AuthObsClients map[string][]*obsclient.AuthObsClient
}

func (n *RPCHandles) RndEthClient() ethadapter.EthClient {
	return n.EthClients[rand.Intn(len(n.EthClients))] //nolint:gosec
}

// ObscuroWalletRndClient fetches an RPC client connected to a random L2 node for a given wallet
func (n *RPCHandles) ObscuroWalletRndClient(wallet wallet.Wallet) *obsclient.AuthObsClient {
	addr := wallet.Address().String()
	clients := n.AuthObsClients[addr]
	return clients[rand.Intn(len(clients))] //nolint:gosec
}

// ObscuroWalletClient fetches a client for a given wallet address, for a specific node
func (n *RPCHandles) ObscuroWalletClient(walletAddress common.Address, nodeIdx int) *obsclient.AuthObsClient {
	clients := n.AuthObsClients[walletAddress.String()]
	return clients[nodeIdx]
}
