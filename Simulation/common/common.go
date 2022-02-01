package common

type NodeId uint64

type StatsCollector interface {
	L1Reorg(id NodeId)
	L2Recalc(id NodeId)
	NewBlock(block Block)
	NewRollup(rollup Rollup)
}
