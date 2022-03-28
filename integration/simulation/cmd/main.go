package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

// DefaultAverageLatencyToBlockRatio is relative to the block time
// Average eth Block duration=12s, and average eth block latency = 1s
// Determines the broadcast powTime. The lower, the more powTime.
const DefaultAverageLatencyToBlockRatio = 12

// DefaultAverageGossipPeriodToBlockRatio - how long to wait for gossip in L2.
const DefaultAverageGossipPeriodToBlockRatio = 3

func main() {
	//f, err := os.Create("cpu.prof")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//defer f.Close() // error handling omitted for example
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	f1, err := os.Create("simulation_result.txt")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	log.SetLog(f1)

	// define core test parameters
	numberOfNodes := 10
	numberOfWallets := 5
	simulationTimeSecs := 15
	avgBlockDurationUSecs := uint64(25_000)
	avgLatency := avgBlockDurationUSecs / DefaultAverageLatencyToBlockRatio
	avgGossipPeriod := avgBlockDurationUSecs / DefaultAverageGossipPeriodToBlockRatio

	// converted to Us
	simulationTimeUSecs := simulationTimeSecs * 1000 * 1000

	// define network params
	stats := simulation.NewStats(numberOfNodes)

	mockEthNodes, obscuroInMemNodes := simulation.CreateBasicNetworkOfInMemoryNodes(numberOfNodes, avgGossipPeriod, avgBlockDurationUSecs, avgLatency, stats)

	txInjector := simulation.NewTransactionInjector(numberOfWallets, avgBlockDurationUSecs, stats, simulationTimeUSecs, mockEthNodes, obscuroInMemNodes)

	sim := simulation.Simulation{
		MockEthNodes:       mockEthNodes,      // the list of mock ethereum nodes
		ObscuroNodes:       obscuroInMemNodes, //  the list of in memory obscuro nodes
		AvgBlockDuration:   avgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: simulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	sim.Start()
	fmt.Printf("%s\n", simulation.NewOutputStats(&sim))
	sim.Stop()
}
