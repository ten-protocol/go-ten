package simulation

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/google/uuid"
)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, netw network.Network, params params.SimParams) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	stats := stats2.NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	ethClients, obscuroNodes, p2pAddrs := netw.Create(params, stats)

	txInjector := NewTransactionInjector(params.NumberOfWallets, params.AvgBlockDuration, stats, ethClients, obscuroNodes)

	simulation := Simulation{
		EthClients:       ethClients,
		ObscuroNodes:     obscuroNodes,
		ObscuroP2PAddrs:  p2pAddrs,
		AvgBlockDuration: uint64(params.AvgBlockDuration),
		TxInjector:       txInjector,
		SimulationTime:   params.SimulationTime,
		Stats:            stats,
		Params:           &params,
	}

	waitForNodesReady(obscuroNodes)

	// execute the simulation
	simulation.Start()

	// run tests
	checkNetworkValidity(t, &simulation)

	simulation.Stop()

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
	netw.TearDown()
}

func waitForNodesReady(obsNodes []*host.Node) {
	now := time.Now()
	for _, n := range obsNodes {
		for {
			if n.IsReady() {
				fmt.Printf("Node %d is Ready after %s\n", obscurocommon.ShortAddress(n.ID), time.Since(now))
				break
			}
			fmt.Printf("Waiting on Node %d after %s \n", obscurocommon.ShortAddress(n.ID), time.Since(now))
			time.Sleep(time.Millisecond)
		}
	}

}
