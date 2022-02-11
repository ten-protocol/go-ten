package simulation

import (
	"simulation/common"
	common2 "simulation/obscuro/common"
	"sync"
)

// Stats - collects information during the simulation. It can be checked programmatically.
type Stats struct {
	nrMiners         int
	simulationTime   int
	avgBlockDuration uint64
	avgLatency       uint64
	gossipPeriod     uint64

	l1Height      int
	totalL1Blocks int

	l2Height           int
	totalL2Blocks      int
	l2Head             *common2.Rollup
	maxRollupsPerBlock uint32
	nrEmptyBlocks      int

	totalL2Txs  int
	noL1Reorgs  map[common.NodeId]int
	noL2Recalcs map[common.NodeId]int
	// todo - actual avg block Duration

	totalDepositedAmount      uint64
	totalWithdrawnAmount      uint64
	rollupWithMoreRecentProof uint64
	nrTransferTransactions    int
	statsMu                   *sync.RWMutex
}

func NewStats(nrMiners int, simulationTime int, avgBlockDuration uint64, avgLatency uint64, gossipPeriod uint64) Stats {
	return Stats{
		nrMiners:         nrMiners,
		simulationTime:   simulationTime,
		avgBlockDuration: avgBlockDuration,
		avgLatency:       avgLatency,
		gossipPeriod:     gossipPeriod,
		noL1Reorgs:       map[common.NodeId]int{},
		noL2Recalcs:      map[common.NodeId]int{},
		statsMu:          &sync.RWMutex{},
	}
}

func (s *Stats) L1Reorg(id common.NodeId) {
	s.statsMu.Lock()
	s.noL1Reorgs[id]++
	s.statsMu.Unlock()
}

func (s *Stats) L2Recalc(id common.NodeId) {
	s.statsMu.Lock()
	s.noL2Recalcs[id]++
	s.statsMu.Unlock()
}

func (s *Stats) NewBlock(b *common.Block) {
	s.statsMu.Lock()
	//s.l1Height = common.MaxInt(s.l1Height, b.Height)
	s.totalL1Blocks++
	s.maxRollupsPerBlock = common.MaxInt(s.maxRollupsPerBlock, uint32(len(b.Transactions)))
	if len(b.Transactions) == 0 {
		s.nrEmptyBlocks++
	}
	s.statsMu.Unlock()
}

func (s *Stats) NewRollup(r *common2.Rollup) {
	s.statsMu.Lock()
	s.l2Head = r
	s.totalL2Blocks++
	s.totalL2Txs += len(r.Transactions)
	s.statsMu.Unlock()
}

func (s *Stats) Deposit(v uint64) {
	s.statsMu.Lock()
	s.totalDepositedAmount += v
	s.statsMu.Unlock()
}

func (s *Stats) Transfer() {
	s.statsMu.Lock()
	s.nrTransferTransactions++
	s.statsMu.Unlock()
}

func (s *Stats) Withdrawal(v uint64) {
	s.statsMu.Lock()
	s.totalWithdrawnAmount += v
	s.statsMu.Unlock()
}

func (s *Stats) RollupWithMoreRecentProof() {
	s.statsMu.Lock()
	s.rollupWithMoreRecentProof++
	s.statsMu.Unlock()
}
