package simulation

import (
	"math/rand"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/go/log"
	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/google/uuid"
)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, netw network.Network, params *params.SimParams) {
	defer func() {
		// wait until clean up is complete before we log the lingering goroutine count
		log.Info("goroutine leak monitor - simulation end - %d goroutines currently running", runtime.NumGoroutine())
	}()
	log.Info("goroutine leak monitor - simulation start - %d goroutines currently running", runtime.NumGoroutine())
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	stats := stats2.NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	ethClients, obscuroClients, p2pAddrs := netw.Create(params, stats)
	defer netw.TearDown()

	txInjector := NewTransactionInjector(params.NumberOfObscuroWallets, params.AvgBlockDuration, stats, ethClients, obscuroClients)

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

	// wait for p2p addresses to be connectable (fudge because we don't handle dropped messages)
	for simulation.Params.WaitForP2PConnections && !allP2PAddressesReady(simulation.ObscuroP2PAddrs) {
		time.Sleep(simulation.Params.AvgBlockDuration * 10)
		log.Info("Waiting for P2P connections to be available.")
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkNetworkValidity(t, &simulation)

	simulation.Stop()

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
}

func allP2PAddressesReady(addrs []string) bool {
	for _, a := range addrs {
		conn, err := net.Dial("tcp", a)
		if err != nil {
			return false
		}
		// we don't worry about failure while closing, it connected successfully so let test proceed
		_ = conn.Close()
	}
	return true
}
