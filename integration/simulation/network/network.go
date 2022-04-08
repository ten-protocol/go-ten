package network

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

// Network is responsible with knowing how to manage the lifecycle of networks of Ethereum or Obscuro nodes.
// These networks can be composed of in-memory go-routines or of fully fledged existing nodes like Ropsten.
// Implementation notes:
// - This is a work in progress, so there is a lot of code duplication in the implementations
// - Once we implement a few more versions: for example using Ganache, or using enclaves running in azure, etc, we'll revisit and create better abstractions.
type Network interface {
	// Create - returns a group of started Ethereum nodes, a group of started Obscuro nodes, and the Obscuro nodes' P2P addresses.
	// todo - return interfaces to RPC handles to the nodes
	Create(params params.SimParams, stats *stats.Stats) ([]*ethereum_mock.Node, []*host.Node, []string)
	TearDown()
}
