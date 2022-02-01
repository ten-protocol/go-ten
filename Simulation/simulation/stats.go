package simulation

import (
	"simulation/common"
	"sync"
)

// Stats - collects information during the simulation. It can be checked programmatically.
type Stats struct {
	nrMiners         int
	simulationTime   int
	avgBlockDuration uint64
	avgLatency       uint64
	gossipPeriod     uint64

	l1Height      uint32
	totalL1Blocks int

	l2Height           uint32
	totalL2Blocks      int
	l2Head             *common.Rollup
	maxRollupsPerBlock uint32
	nrEmptyBlocks      int

	totalL2Txs  int
	noL1Reorgs  map[common.NodeId]int
	noL2Recalcs map[common.NodeId]int
	// todo - actual avg block Duration

	totalDepositedAmount   uint64
	totalWithdrawnAmount   uint64
	nrTransferTransactions int
}

var statsMu = &sync.RWMutex{}

func NewStats(nrMiners int, simulationTime int, avgBlockDuration uint64, avgLatency uint64, gossipPeriod uint64) Stats {
	return Stats{
		nrMiners:         nrMiners,
		simulationTime:   simulationTime,
		avgBlockDuration: avgBlockDuration,
		avgLatency:       avgLatency,
		gossipPeriod:     gossipPeriod,
		noL1Reorgs:       map[common.NodeId]int{},
		noL2Recalcs:      map[common.NodeId]int{},
	}
}

func (s *Stats) L1Reorg(id common.NodeId) {
	statsMu.Lock()
	s.noL1Reorgs[id]++
	statsMu.Unlock()
}

func (s *Stats) L2Recalc(id common.NodeId) {
	statsMu.Lock()
	s.noL2Recalcs[id]++
	statsMu.Unlock()
}

func (s *Stats) NewBlock(b common.Block) {
	statsMu.Lock()
	s.l1Height = common.MaxInt(s.l1Height, b.Height())
	s.totalL1Blocks++
	s.maxRollupsPerBlock = common.MaxInt(s.maxRollupsPerBlock, uint32(len(b.Txs())))
	if len(b.Txs()) == 0 {
		s.nrEmptyBlocks++
	}
	statsMu.Unlock()
}

func (s *Stats) NewRollup(r common.Rollup) {
	statsMu.Lock()
	s.l2Height = common.MaxInt(s.l2Height, r.Height())
	s.l2Head = &r
	s.totalL2Blocks++
	s.totalL2Txs += len(r.L2Txs())
	statsMu.Unlock()
}

func (s *Stats) Deposit(v uint64) {
	statsMu.Lock()
	s.totalDepositedAmount += v
	statsMu.Unlock()
}

func (s *Stats) Transfer() {
	statsMu.Lock()
	s.nrTransferTransactions++
	statsMu.Unlock()
}

func (s *Stats) Withdrawal(v uint64) {
	statsMu.Lock()
	s.totalWithdrawnAmount += v
	statsMu.Unlock()
}
