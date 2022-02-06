package common

type NodeId uint64

// todo - use proper crypto
//type Address = uuid.UUID
type Address = uint32

type StatsCollector interface {
	L1Reorg(id NodeId)
	L2Recalc(id NodeId)
	NewBlock(block Block)
	NewRollup(rollup Rollup)
}
