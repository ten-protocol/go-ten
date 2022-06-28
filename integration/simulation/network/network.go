package network

import (
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
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
	Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []rpcclientlib.Client, error)
	TearDown()
}
