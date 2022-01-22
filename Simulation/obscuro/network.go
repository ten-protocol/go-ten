package obscuro

import (
	"sync"
)

// NetworkCfg - models a full network including artificial random latencies
type NetworkCfg struct {
	allMiners []L1Miner
	allAgg    []L2Agg
	delay     Latency // the latency
	Stats     *Stats
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
	c.Stats.l1Height = Max(c.Stats.l1Height, b.height)
	c.Stats.totalL1++
	c.Stats.maxRollupsPerBlock = Max(c.Stats.maxRollupsPerBlock, len(b.txs))
	if len(b.txs) == 0 {
		c.Stats.nrEmptyBlocks++
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
		d := Max(c.delay()/2, 1)
		Schedule(d, func() { t.L1P2PGossipTx(tx) })
	}

	// collect Stats
	if tx.txType == RollupTx {
		statsMu.Lock()
		c.Stats.l2Height = Max(c.Stats.l2Height, tx.rollup.height)
		c.Stats.l2Head = tx.rollup
		c.Stats.totalL2++
		c.Stats.totalL2Txs += len(tx.rollup.txs)
		statsMu.Unlock()
	}
}

func (c *NetworkCfg) broadcastL2Tx(tx *L2Tx) {
	for _, a := range c.allAgg {
		t := a
		Schedule(c.delay()/2, func() { t.L2P2PReceiveTx(tx) })
	}
}
