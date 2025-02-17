package simulation

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ethereum/go-ethereum"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"

	simstats "github.com/ten-protocol/go-ten/integration/simulation/stats"

	"github.com/google/uuid"
)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, netw network.Network, params *params.SimParams) {
	defer func() {
		// wait until clean up is complete before we log the lingering goroutine count
		testlog.Logger().Info(fmt.Sprintf("goroutine leak monitor - simulation end - %d goroutines currently running", runtime.NumGoroutine()))
	}()
	testlog.Logger().Info(fmt.Sprintf("goroutine leak monitor - simulation start - %d goroutines currently running", runtime.NumGoroutine()))
	rand.Seed(time.Now().UnixNano()) //nolint: staticcheck
	uuid.EnableRandPool()

	stats := simstats.NewStats(params.NumberOfNodes)

	fmt.Printf("Creating network\n")
	testlog.Logger().Info("Creating network")
	defer netw.TearDown()
	networkClients, err := netw.Create(params, stats)
	// Return early if the network was not created
	if err != nil {
		fmt.Printf("Could not run test: %s\n", err)
		t.Fatal()
	}

	txInjector := NewTransactionInjector(
		params.AvgBlockDuration,
		stats,
		networkClients,
		params.Wallets,
		&params.L1TenData.MgmtContractAddress,
		params.MgmtContractLib,
		params.ERC20ContractLib,
		0,
		params,
	)

	simulation := Simulation{
		RPCHandles:       networkClients,
		AvgBlockDuration: uint64(params.AvgBlockDuration),
		TxInjector:       txInjector,
		SimulationTime:   params.SimulationTime,
		Stats:            stats,
		Params:           params,
		LogChannels:      make(map[string][]chan types.Log),
		Subscriptions:    []ethereum.Subscription{},
	}

	// execute the simulation
	fmt.Printf("Starting simulation\n")
	testlog.Logger().Info("Starting simulation")
	simulation.Start()

	// run tests
	fmt.Printf("Validating simulation results\n")
	testlog.Logger().Info("Validating simulation results")
	checkNetworkValidity(t, &simulation)

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
	testlog.Logger().Info(fmt.Sprintf("Simulation results:%+v", NewOutputStats(&simulation)))

	fmt.Printf("Stopping simulation\n")
	testlog.Logger().Info("Stopping simulation")
	simulation.Stop()
}
