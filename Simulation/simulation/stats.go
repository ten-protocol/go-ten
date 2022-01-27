package simulation

import (
	"simulation/common"
	"sync"
)

// Stats - collects information during the simulation. It can be checked programmatically.
type Stats struct {
	nrMiners         int
	simulationTime   int
	avgBlockDuration int
	avgLatency       int
	gossipPeriod     int

	l1Height      int
	totalL1Blocks int

	l2Height           int
	totalL2Blocks      int
	l2Head             *common.Rollup
	maxRollupsPerBlock int
	nrEmptyBlocks      int

	totalL2Txs  int
	noL1Reorgs  map[common.NodeId]int
	noL2Recalcs map[common.NodeId]int
	// todo - actual avg block Duration

	totalDepositedAmount   int
	totalWithdrawnAmount   int
	nrTransferTransactions int
}

var statsMu = &sync.RWMutex{}

func NewStats(nrMiners int, simulationTime int, avgBlockDuration int, avgLatency int, gossipPeriod int) Stats {
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
	s.l1Height = common.Max(s.l1Height, b.Height())
	s.totalL1Blocks++
	s.maxRollupsPerBlock = common.Max(s.maxRollupsPerBlock, len(b.Txs()))
	if len(b.Txs()) == 0 {
		s.nrEmptyBlocks++
	}
	statsMu.Unlock()
}

func (s *Stats) NewRollup(r common.Rollup) {
	statsMu.Lock()
	s.l2Height = common.Max(s.l2Height, r.Height())
	s.l2Head = &r
	s.totalL2Blocks++
	s.totalL2Txs += len(r.L2Txs())
	statsMu.Unlock()
}

func (s *Stats) Deposit(v int) {
	statsMu.Lock()
	s.totalDepositedAmount += v
	statsMu.Unlock()
}

func (s *Stats) Transfer() {
	statsMu.Lock()
	s.nrTransferTransactions++
	statsMu.Unlock()
}

func (s *Stats) Withdrawal(v int) {
	statsMu.Lock()
	s.totalWithdrawnAmount += v
	statsMu.Unlock()
}
