package main

// DefaultAverageLatencyToBlockRatio is relative to the block time
// Average eth Block duration=12s, and average eth block latency = 1s
// Determines the broadcast powTime. The lower, the more powTime.
const DefaultAverageLatencyToBlockRatio = 12

// DefaultAverageGossipPeriodToBlockRatio - how long to wait for gossip in L2.
const DefaultAverageGossipPeriodToBlockRatio = 3

func main() {
	////f, err := os.Create("cpu.prof")
	////if err != nil {
	////	log.Fatal("could not create CPU profile: ", err)
	////}
	////defer f.Close() // error handling omitted for example
	////if err := pprof.StartCPUProfile(f); err != nil {
	////	log.Fatal("could not start CPU profile: ", err)
	////}
	////defer pprof.StopCPUProfile()
	//rand.Seed(time.Now().UnixNano())
	//uuid.EnableRandPool()
	//
	//f1, err := os.Create("simulation_result.txt")
	//if err != nil {
	//	panic(err)
	//}
	//defer f1.Close()
	//log.SetLog(f1)
	//
	//// define core test parameters
	//numberOfNodes := 10
	//simulationTime := 15
	//avgBlockDuration := uint64(25_000)
	//avgLatency := avgBlockDuration / DefaultAverageLatencyToBlockRatio
	//avgGossipPeriod := avgBlockDuration / DefaultAverageGossipPeriodToBlockRatio
	//
	//// define network params
	//stats := simulation.NewStats(numberOfNodes)
	//l1NetworkConfig := simulation.NewL1Network(avgBlockDuration, avgLatency, stats)
	//l2NetworkCfg := simulation.NewL2Network(avgBlockDuration, avgLatency)
	//
	//// define instances of the simulation mechanisms
	//txManager := simulation.NewTransactionManager(5, l1NetworkConfig, l2NetworkCfg, avgBlockDuration, stats)
	//sim := simulation.NewSimulation(
	//	numberOfNodes,
	//	l1NetworkConfig,
	//	l2NetworkCfg,
	//	avgBlockDuration,
	//	avgGossipPeriod,
	//	stats,
	//)
	//
	//// execute the simulation
	//sim.Start(txManager, simulationTime)
	//fmt.Printf("%s\n", simulation.NewOutputStats(sim))
}
