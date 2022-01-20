package main

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"simulation/obscuro"
	"time"
)

// DEFAULT_AVERAGE_LATENCY_TO_BLOCK_RATIO is relative to the block time
// Average eth Block duration=12s, and average eth block latency = 1s
// Determines the broadcast powTime. The lower, the more powTime.
const DEFAULT_AVERAGE_LATENCY_TO_BLOCK_RATIO = 12

// DEFAULT_AVERAGE_GOSSIP_PERIOD_TO_BLOCK_RATIO - how long to wait for gossip in L2.
const DEFAULT_AVERAGE_GOSSIP_PERIOD_TO_BLOCK_RATIO = 3

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

	f, err := os.Create("simulation_result.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	obscuro.SetLog(f)

	blockDuration := 10_000
	stats := obscuro.RunSimulation(10, 10, 30, blockDuration, blockDuration/12, blockDuration/2)
	fmt.Printf("%#v\n", stats)
}
