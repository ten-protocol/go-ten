package simulation

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// OutputStats decouples the processing of data and the collection of statistics
// there's a bit more to do around this, this serves as a first iteration
type OutputStats struct {
	simulation *Simulation

	l2RollupCountInHeaders    int // Number of rollups counted while node rollup header traversing
	l2RollupCountInL1Blocks   int // Number of rollups counted while traversing the node block header and searching the txs
	l2RollupTxCountInL1Blocks int // Number of rollup Txs counted while traversing the node block header
	l1Height                  int // Last known l1 block height
	l2Height                  int // Last known l2 block height
}

// NewOutputStats processes the simulation and retrieves the output statistics
func NewOutputStats(simulation *Simulation) *OutputStats {
	outputStats := &OutputStats{
		simulation: simulation,
	}

	outputStats.countRollups()
	outputStats.populateHeights()

	return outputStats
}

func (o *OutputStats) populateHeights() {
	o.l1Height = int(o.simulation.ObscuroNodes[0].DB().GetCurrentBlockHead().Height)
	o.l2Height = int(o.simulation.ObscuroNodes[0].DB().GetCurrentRollupHead().Height)
}

func (o *OutputStats) countRollups() {
	l1Node := o.simulation.EthClients[0]
	l2Node := o.simulation.ObscuroNodes[0]

	// iterate the Node Headers and get the rollups
	for header := l2Node.DB().GetCurrentRollupHead(); header != nil && header.ID != obscurocommon.GenesisHash; header = l2Node.DB().GetRollupHeader(header.Parent) {
		o.l2RollupCountInHeaders++
	}

	// iterate the L1 Blocks and get the rollups
	for header := l2Node.DB().GetCurrentBlockHead(); header != nil && header.ID != obscurocommon.GenesisHash; header = l2Node.DB().GetBlockHeader(header.Parent) {
		block, found := l1Node.FetchBlock(header.ID)
		if !found {
			panic("expected l1 block not found")
		}
		for _, tx := range block.Transactions() {
			txData := obscurocommon.TxData(tx)
			if txData.TxType == obscurocommon.RollupTx {
				r := nodecommon.DecodeRollupOrPanic(txData.Rollup)
				if l1Node.IsBlockAncestor(*block, r.Header.L1Proof) {
					o.l2RollupCountInL1Blocks++
					o.l2RollupTxCountInL1Blocks += len(r.Transactions)
				}
			}
		}
	}
}

func (o *OutputStats) String() string {
	return fmt.Sprintf("\n"+
		"nrMiners: %d\n"+
		"l1Height: %d\n"+
		"l2Height: %d\n"+
		"totalL1Blocks: %d\n"+
		"totalL2Blocks: %v\n"+
		"l2RollupCountInHeaders: %d\n"+
		"l2RollupCountInL1Blocks: %d\n"+
		"l2RollupTxCountInL1Blocks: %d\n"+
		"maxRollupsPerBlock: %d \n"+
		"nrEmptyBlocks: %d\n"+
		"noL1Reorgs: %+v\n"+
		"noL2Recalcs: %+v\n"+
		"totalDepositedAmount: %d\n"+
		"totalWithdrawnAmount: %d\n"+
		"rollupWithMoreRecentProof: %d\n"+
		"nrTransferTransactions: %d\n",
		o.simulation.Stats.NrMiners,
		o.l1Height,
		o.l2Height,
		o.simulation.Stats.TotalL1Blocks,
		o.simulation.Stats.NoL2Blocks,
		o.l2RollupCountInHeaders,
		o.l2RollupCountInL1Blocks,
		o.l2RollupTxCountInL1Blocks,
		o.simulation.Stats.MaxRollupsPerBlock,
		o.simulation.Stats.NrEmptyBlocks,
		o.simulation.Stats.NoL1Reorgs,
		o.simulation.Stats.NoL2Recalcs,
		o.simulation.Stats.TotalDepositedAmount,
		o.simulation.Stats.TotalWithdrawalRequestedAmount,
		o.simulation.Stats.RollupWithMoreRecentProofCount,
		o.simulation.Stats.NrTransferTransactions,
	)
}
