package network

import (
	"github.com/ethereum/go-ethereum/common"
	"math/rand"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

// Network is responsible with knowing how to manage the lifecycle of networks of Ethereum or Obscuro nodes.
// These networks can be composed of in-memory go-routines or of fully fledged existing nodes like Ropsten.
// Implementation notes:
// - This is a work in progress, so there is a lot of code duplication in the implementations
// - Once we implement a few more versions: for example using Geth, or using enclaves running in azure, etc, we'll revisit and create better abstractions.
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
	ObscuroClients []rpcclientlib.Client

	// an RPC client per node per wallet, with a viewing key setup (on the client and registered on its corresponding host enclave),
	//	to mimic user acc interaction via a wallet extension
	// map of address hex string -> slice of clients (indexed by node)
	VirtualWalletExtensionClients map[string][]rpcclientlib.Client
}

func (n *RPCHandles) RndEthClient() ethadapter.EthClient {
	return n.EthClients[rand.Intn(len(n.EthClients))] //nolint:gosec
}

func (n *RPCHandles) RndObscuroClient() rpcclientlib.Client {
	return n.ObscuroClients[rand.Intn(len(n.ObscuroClients))] //nolint:gosec
}

// ObscuroWalletRndClient fetches a random client for given wallet
func (n *RPCHandles) ObscuroWalletRndClient(wallet wallet.Wallet) rpcclientlib.Client {
	addr := wallet.Address().String()
	clients := n.VirtualWalletExtensionClients[addr]
	return clients[rand.Intn(len(clients))]
}

// ObscuroWalletClient fetches a client for a given wallet address, for a specific node
func (n *RPCHandles) ObscuroWalletClient(walletAddress common.Address, nodeIdx int) rpcclientlib.Client {
	clients := n.VirtualWalletExtensionClients[walletAddress.String()]
	return clients[nodeIdx]
}
