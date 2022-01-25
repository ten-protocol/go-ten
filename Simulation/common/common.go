package common

type NodeId int

type L1Network interface {
	BroadcastBlockL1(b Block)
	BroadcastL1Tx(tx L1Tx)
}

type L2Network interface {
	BroadcastRollupL2(r Rollup)
	BroadcastL2Tx(tx L2Tx)
}

type NotifyNewBlock interface {
	RPCNewHead(b Block)
}

type StatsCollector interface {
	L1Reorg(id NodeId)
	NewBlock(block Block)
	NewRollup(rollup Rollup)
}
