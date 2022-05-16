package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/google/uuid"
)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, netw network.Network, params *params.SimParams) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	stats := stats2.NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	ethClients, obscuroClients, p2pAddrs := netw.Create(params, stats)

	var workerWallet wallet.Wallet
	if len(params.EthWallets) > 0 {
		workerWallet = params.EthWallets[len(params.EthWallets)-1]
	}

	txInjector := NewTransactionInjector(
		params.NumberOfObscuroWallets,
		params.AvgBlockDuration,
		stats,
		ethClients,
		workerWallet,
		params.MgmtContractAddr,
		params.ERC20ContractAddr,
		obscuroClients,
	)

	simulation := Simulation{
		EthClients:       ethClients,
		ObscuroClients:   obscuroClients,
		ObscuroP2PAddrs:  p2pAddrs,
		AvgBlockDuration: uint64(params.AvgBlockDuration),
		TxInjector:       txInjector,
		SimulationTime:   params.SimulationTime,
		Stats:            stats,
		Params:           params,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkNetworkValidity(t, &simulation)

	simulation.Stop()

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
	netw.TearDown()
}
