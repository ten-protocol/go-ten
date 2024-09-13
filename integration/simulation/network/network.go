package network

import (
	"math/rand"
	"testing"

	"github.com/ten-protocol/go-ten/go/rpc"

	"github.com/ten-protocol/go-ten/go/obsclient"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/integration/simulation/stats"
)

// Network is responsible with knowing how to manage the lifecycle of networks of Ethereum or TEN nodes.
// These networks can be composed of in-memory go-routines or of fully fledged existing nodes like Ropsten.
// Implementation notes:
// - This is a work in progress, so there is a lot of code duplication in the implementations
// - Once we implement a few more versions: for example using Geth, we'll revisit and create better abstractions.
type Network interface {
	// Create - returns the started Ethereum nodes and the started TEN node clients.
	// Responsible with spinning up all resources required for the test
	// Return an error in case it cannot start for an expected reason, otherwise it panics.
	Create(params *params.SimParams, stats *stats.Stats, t *testing.T) (*RPCHandles, error)
	TearDown()
}

type RPCHandles struct {
	// an eth client per eth node in the network
	EthClients []ethadapter.EthClient

	// A TEN client per TEN node in the network.
	TenClients []*obsclient.ObsClient
	// An RPC client per TEN node in the network (used for APIs that don't have methods on `ObsClient`.
	RPCClients []rpc.Client

	// an RPC client per node per wallet, with a viewing key set up (on the client and registered on its corresponding host enclave),
	//	to mimic user acc interaction via a wallet extension
	// map of owner addresses to RPC clients for that owner (one per L2 node)
	// todo (@matt) - simplify this with a client per node when we have clients that can support multiple wallets
	AuthObsClients map[string][]*obsclient.AuthObsClient
}

func (n *RPCHandles) RndEthClient() ethadapter.EthClient {
	return n.EthClients[rand.Intn(len(n.EthClients))] //nolint:gosec
}

// TenWalletRndClient fetches an RPC client connected to a random L2 node for a given wallet
func (n *RPCHandles) TenWalletRndClient(wallet wallet.Wallet) *obsclient.AuthObsClient {
	addr := wallet.Address().String()
	clients := n.AuthObsClients[addr]
	return clients[rand.Intn(len(clients))] //nolint:gosec
}

// TenWalletClient fetches a client for a given wallet address, for a specific node
func (n *RPCHandles) TenWalletClient(walletAddress common.Address, nodeIdx int) *obsclient.AuthObsClient {
	clients := n.AuthObsClients[walletAddress.String()]
	return clients[nodeIdx]
}
