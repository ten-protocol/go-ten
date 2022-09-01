package stats

import (
	"math/big"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"
)

// Stats - collects information during the simulation. It can be checked programmatically.
// Todo - this is a temporary placeholder until we introduce a proper metrics framework like prometheus
type Stats struct {
	NrMiners int

	TotalL1Blocks uint64

	MaxRollupsPerBlock uint32
	NrEmptyBlocks      int

	NoL1Reorgs  map[gethcommon.Address]int
	NoL2Recalcs map[gethcommon.Address]int
	NoL2Blocks  map[gethcommon.Address]uint64
	// todo - actual avg block Duration

	TotalDepositedAmount           *big.Int
	TotalWithdrawalRequestedAmount *big.Int
	RollupWithMoreRecentProofCount uint64
	NrTransferTransactions         int
	statsMu                        *sync.RWMutex
}

func NewStats(nrMiners int) *Stats {
	return &Stats{
		NrMiners:                       nrMiners,
		NoL1Reorgs:                     map[gethcommon.Address]int{},
		NoL2Recalcs:                    map[gethcommon.Address]int{},
		NoL2Blocks:                     map[gethcommon.Address]uint64{},
		TotalDepositedAmount:           big.NewInt(0),
		TotalWithdrawalRequestedAmount: big.NewInt(0),
		statsMu:                        &sync.RWMutex{},
	}
}

func (s *Stats) L1Reorg(id gethcommon.Address) {
	s.statsMu.Lock()
	s.NoL1Reorgs[id]++
	s.statsMu.Unlock()
}

func (s *Stats) L2Recalc(id gethcommon.Address) {
	s.statsMu.Lock()
	s.NoL2Recalcs[id]++
	s.statsMu.Unlock()
}

func (s *Stats) NewBlock(b *types.Block) {
	s.statsMu.Lock()
	// s.l1Height = nodecommon.MaxInt(s.l1Height, b.Number)
	s.MaxRollupsPerBlock = common.MaxInt(s.MaxRollupsPerBlock, uint32(len(b.Transactions())))
	if len(b.Transactions()) == 0 {
		s.NrEmptyBlocks++
	}
	s.statsMu.Unlock()
}

func (s *Stats) NewRollup(node gethcommon.Address) {
	s.statsMu.Lock()
	s.NoL2Blocks[node]++
	s.statsMu.Unlock()
}

func (s *Stats) Deposit(v *big.Int) {
	s.statsMu.Lock()
	s.TotalDepositedAmount.Add(s.TotalDepositedAmount, v)
	s.statsMu.Unlock()
}

func (s *Stats) Transfer() {
	s.statsMu.Lock()
	s.NrTransferTransactions++
	s.statsMu.Unlock()
}

func (s *Stats) Withdrawal(v *big.Int) {
	s.statsMu.Lock()
	s.TotalWithdrawalRequestedAmount = s.TotalWithdrawalRequestedAmount.Add(s.TotalWithdrawalRequestedAmount, v)
	s.statsMu.Unlock()
}

func (s *Stats) RollupWithMoreRecentProof() {
	s.statsMu.Lock()
	s.RollupWithMoreRecentProofCount++
	s.statsMu.Unlock()
}
