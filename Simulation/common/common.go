package common

type NodeId int

type StatsCollector interface {
	L1Reorg(id NodeId)
	L2Recalc(id NodeId)
	NewBlock(block Block)
	NewRollup(rollup Rollup)
}
