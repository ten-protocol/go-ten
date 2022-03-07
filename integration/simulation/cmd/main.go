package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

// DefaultAverageLatencyToBlockRatio is relative to the block time
// Average eth Block duration=12s, and average eth block latency = 1s
// Determines the broadcast powTime. The lower, the more powTime.
const DefaultAverageLatencyToBlockRatio = uint64(12)

// DefaultAverageGossipPeriodToBlockRatio - how long to wait for gossip in L2.
const DefaultAverageGossipPeriodToBlockRatio = uint64(3)

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
	simulationTime := 15
	avgBlockDuration := DefaultAverageLatencyToBlockRatio
	avgLatency := avgBlockDuration / 15
	avgGossipPeriod := DefaultAverageGossipPeriodToBlockRatio

	// define network params
	stats := simulation.NewStats(numberOfNodes)
	l1NetworkConfig := simulation.NewL1Network(avgLatency, stats)
	l2NetworkCfg := simulation.NewL2Network(avgLatency)

	// define instances of the simulation mechanisms
	txManager := simulation.NewTransactionManager(5, l1NetworkConfig, l2NetworkCfg, avgBlockDuration, stats)
	simulationNetwork := simulation.NewSimulationNetwork(
		numberOfNodes,
		l1NetworkConfig,
		l2NetworkCfg,
		avgBlockDuration,
		avgGossipPeriod,
		stats,
	)

	// execute the simulation
	simulation.RunSimulation(txManager, simulationNetwork, simulationTime)
	fmt.Printf("%#v\n", l1NetworkConfig.Stats)
}
