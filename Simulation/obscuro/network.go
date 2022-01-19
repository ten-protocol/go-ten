package obscuro

import (
	"sync"
)

// NetworkCfg - models a full network including artificial random latencies
type NetworkCfg struct {
	allMiners []L1Miner
	allAgg    []L2Agg
	delay     Latency // the latency
	stats     *Stats
}

var statsMu = &sync.RWMutex{}

// broadcast a block to the l1 nodes
func (c *NetworkCfg) broadcastBlockL1(b *Block) {
	for _, m := range c.allMiners {
		if m.id != b.miner.id {
			t := m
			Schedule(c.delay(), func() { t.P2PReceiveBlock(b) })
		}
	}
	statsMu.Lock()
	c.stats.l1Height = Max(c.stats.l1Height, b.height)
	c.stats.totalL1++
	c.stats.maxRollupsPerBlock = Max(c.stats.maxRollupsPerBlock, len(b.txs))
	if len(b.txs) == 0 {
		c.stats.nrEmptyBlocks++
	}
	statsMu.Unlock()
}

// Broadcasts the rollup to all L2 peers
func (c *NetworkCfg) broadcastRollupL2(r *Rollup) {
	for _, a := range c.allAgg {
		if a.id != r.agg.id {
			t := a
			Schedule(c.delay(), func() { t.L2P2PGossipRollup(r) })
		}
	}
	//statsMu.Lock()
	//statsMu.Unlock()
}

// Broadcasts the L1 tx containing the rollup to the L1 network
func (c *NetworkCfg) broadcastL1Tx(tx *L1Tx) {
	for _, m := range c.allMiners {
		t := m
		// the time to broadcast a tx is half that of a L1 block, because it is smaller.
		// todo - find a better way to express this
		Schedule(c.delay()/2, func() { t.L1P2PGossipTx(tx) })
	}

	// collect stats
	if tx.txType == RollupTx {
		statsMu.Lock()
		c.stats.l2Height = Max(c.stats.l2Height, tx.rollup.height)
		c.stats.totalL2++
		c.stats.totalL2Txs += len(tx.rollup.txs)
		statsMu.Unlock()
	}
}

func (c *NetworkCfg) broadcastL2Tx(tx *L2Tx) {
	for _, a := range c.allAgg {
		t := a
		Schedule(c.delay()/2, func() { t.L2P2PReceiveTx(tx) })
	}
}
