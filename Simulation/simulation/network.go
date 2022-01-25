package simulation

import (
	"simulation/common"
	"simulation/ethereum-mock"
	"simulation/obscuro"
	"sync"
)

// NetworkCfg - models a full network including artificial random latencies
type NetworkCfg struct {
	allMiners []ethereum_mock.L1Miner
	allAgg    []obscuro.L2Agg
	delay     common.Latency // the latency
	Stats     *Stats
}

var statsMu = &sync.RWMutex{}

// broadcast a block to the l1 nodes
func (c NetworkCfg) BroadcastBlockL1(b common.Block) {
	for _, m := range c.allMiners {
		if m.Id != b.Miner {
			t := m
			common.Schedule(c.delay(), func() { t.P2PReceiveBlock(&b) })
		}
	}
	statsMu.Lock()
	c.Stats.l1Height = common.Max(c.Stats.l1Height, b.Height())
	c.Stats.totalL1++
	c.Stats.maxRollupsPerBlock = common.Max(c.Stats.maxRollupsPerBlock, len(b.Txs()))
	if len(b.Txs()) == 0 {
		c.Stats.nrEmptyBlocks++
	}
	statsMu.Unlock()
}

// Broadcasts the rollup to all L2 peers
func (c NetworkCfg) BroadcastRollupL2(r common.Rollup) {
	for _, a := range c.allAgg {
		if a.Id != r.Agg {
			t := a
			common.Schedule(c.delay(), func() { t.L2P2PGossipRollup(&r) })
		}
	}
	//statsMu.Lock()
	//statsMu.Unlock()
}

// Broadcasts the L1 tx containing the rollup to the L1 network
func (c NetworkCfg) BroadcastL1Tx(tx common.L1Tx) {
	for _, m := range c.allMiners {
		t := m
		// the time to broadcast a tx is half that of a L1 block, because it is smaller.
		// todo - find a better way to express this
		d := common.Max(c.delay()/2, 1)
		common.Schedule(d, func() { t.L1P2PGossipTx(&tx) })
	}

	// collect Stats
	if tx.TxType == common.RollupTx {
		statsMu.Lock()
		c.Stats.l2Height = common.Max(c.Stats.l2Height, tx.Rollup.Height())
		c.Stats.l2Head = &tx.Rollup
		c.Stats.totalL2++
		c.Stats.totalL2Txs += len(tx.Rollup.L2Txs())
		statsMu.Unlock()
	}
}

func (c NetworkCfg) BroadcastL2Tx(tx common.L2Tx) {
	for _, a := range c.allAgg {
		t := a
		common.Schedule(c.delay()/2, func() { t.L2P2PReceiveTx(tx) })
	}
}
