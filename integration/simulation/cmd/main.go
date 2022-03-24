package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

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
	simulationTime := 15
	avgBlockDuration := uint64(25_000)
	avgLatency := avgBlockDuration / DefaultAverageLatencyToBlockRatio
	avgGossipPeriod := avgBlockDuration / DefaultAverageGossipPeriodToBlockRatio

	// define network params
	p2pNetwork := p2p.NewP2P
	stats := simulation.NewStats(numberOfNodes)
	l1NetworkConfig := simulation.NewL1Network(avgBlockDuration, avgLatency, stats)
	l2NetworkCfg := simulation.NewL2Network(numberOfNodes, avgBlockDuration, avgLatency, p2pNetwork)

	// define instances of the simulation mechanisms
	txManager := simulation.NewTransactionManager(5, l1NetworkConfig, l2NetworkCfg, avgBlockDuration, stats)
	sim := simulation.NewSimulation(numberOfNodes, l1NetworkConfig, l2NetworkCfg, avgBlockDuration, avgGossipPeriod, false, stats, p2pNetwork)

	// execute the simulation
	sim.Start(txManager, simulationTime)
	fmt.Printf("%s\n", simulation.NewOutputStats(sim))
	sim.Stop()
}
