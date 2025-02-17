package stats

import (
	"math/big"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
)

// Stats - collects information during the simulation. It can be checked programmatically.
type Stats struct {
	NrMiners int

	TotalL1Blocks uint64

	MaxRollupsPerBlock uint32
	NrEmptyBlocks      int

	NoL1Reorgs  map[gethcommon.Address]int
	NoL2Recalcs map[gethcommon.Address]int
	NoL2Blocks  map[int]uint64

	TotalDepositedAmount           *big.Int
	TotalWithdrawalRequestedAmount *big.Int
	RollupWithMoreRecentProofCount uint64
	NrTransferTransactions         int
	NrNativeTransferTransactions   int
	statsMu                        *sync.RWMutex
}

func NewStats(nrMiners int) *Stats {
	return &Stats{
		NrMiners:                       nrMiners,
		NoL1Reorgs:                     map[gethcommon.Address]int{},
		NoL2Recalcs:                    map[gethcommon.Address]int{},
		NoL2Blocks:                     map[int]uint64{},
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

func (s *Stats) NewBlock(b *types.Block) {
	s.statsMu.Lock()
	// s.l1Height = nodecommon.MaxInt(s.l1Height, b.Number)
	s.MaxRollupsPerBlock = common.MaxInt(s.MaxRollupsPerBlock, uint32(len(b.Transactions())))
	if len(b.Transactions()) == 0 {
		s.NrEmptyBlocks++
	}
	s.statsMu.Unlock()
}

func (s *Stats) NewRollup(nodeIdx int) {
	s.statsMu.Lock()
	s.NoL2Blocks[nodeIdx]++
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

func (s *Stats) NativeTransfer() {
	s.statsMu.Lock()
	s.NrNativeTransferTransactions++
	s.statsMu.Unlock()
}

func (s *Stats) Withdrawal(v *big.Int) {
	s.statsMu.Lock()
	s.TotalWithdrawalRequestedAmount = s.TotalWithdrawalRequestedAmount.Add(s.TotalWithdrawalRequestedAmount, v)
	s.statsMu.Unlock()
}
