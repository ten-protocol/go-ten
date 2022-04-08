package params

import "time"

// SimParams are the parameters for setting up the simulation.
type SimParams struct {
	NumberOfNodes   int
	NumberOfWallets int

	// A critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	AvgBlockDurationUSecs uint64
	AvgNetworkLatency     uint64 // artificial latency injected between sending and receiving messages on the mock network
	AvgGossipPeriod       uint64 // POBI protocol setting

	SimulationTime time.Duration // how long the simulations should run for

	// EfficiencyThresholds represents an acceptable "dead blocks" percentage for this simulation.
	// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
	// We test the results against this threshold to catch eventual protocol errors.
	L1EfficiencyThreshold     float64
	L2EfficiencyThreshold     float64
	L2ToL1EfficiencyThreshold float64
}
