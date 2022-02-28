package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
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
	common.SetLog(f1)

	numberOfNodes := 10
	simulationTime := 15
	avgBlockDuration := uint64(20_000)
	avgLatency := avgBlockDuration / 15
	avgGossipPeriod := avgBlockDuration / 3

	stats := simulation.NewStats(numberOfNodes, simulationTime, avgBlockDuration, avgLatency, avgGossipPeriod)

	blockDuration := uint64(25_000)
	l1netw, _ := simulation.RunSimulation(
		simulation.NewTransactionGenerator(5),
		2,
		55,
		blockDuration,
		blockDuration/DefaultAverageLatencyToBlockRatio,
		blockDuration/DefaultAverageGossipPeriodToBlockRatio,
		stats)
	fmt.Printf("%#v\n", l1netw.Stats)
}
