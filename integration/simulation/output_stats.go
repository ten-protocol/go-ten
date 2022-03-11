package simulation

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/common"
	obscuroCommon "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
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

// NewOutputStats process the simulation and retrieves the output statistics
func NewOutputStats(simulation *Simulation) *OutputStats {
	outputStats := &OutputStats{
		simulation: simulation,
	}

	outputStats.countRollups()
	outputStats.populateHeights()

	return outputStats
}

func (o *OutputStats) populateHeights() {
	o.l1Height = int(o.simulation.l2Network.nodes[0].Storage().GetCurrentBlockHead().Height)
	o.l2Height = int(o.simulation.l2Network.nodes[0].Storage().GetCurrentRollupHead().Height)
}

func (o *OutputStats) countRollups() {
	l1Node := o.simulation.l1Network.nodes[0]
	l2Node := o.simulation.l2Network.nodes[0]

	// iterate the Node Headers and get the rollups
	for header := l2Node.Storage().GetCurrentRollupHead(); header != nil; header = l2Node.Storage().GetRollupHeader(header.Parent) {
		o.l2RollupCountInHeaders++
	}

	// iterate the L1 Blocks and get the rollups
	for header := l2Node.Storage().GetCurrentBlockHead(); header != nil; header = l2Node.Storage().GetBlockHeader(header.Parent) {
		block, found := l1Node.Resolver.Resolve(header.ID)
		if !found {
			panic("expected l1 block not found")
		}
		for _, tx := range block.Transactions {
			if tx.TxType == common.RollupTx {
				r := obscuroCommon.DecodeRollup(tx.Rollup)
				if common.IsBlockAncestor(r.Header.L1Proof, block, l1Node.Resolver) {
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
		"totalL2Blocks: %d\n"+
		"l2RollupCountInHeaders: %d\n"+
		"l2RollupCountInL1Blocks: %d\n"+
		"l2RollupTxCountInL1Blocks: %d\n"+
		"maxRollupsPerBlock: %d \n"+
		"nrEmptyBlocks: %d\n"+
		"totalL2Txs: %d\n"+
		"noL1Reorgs: %+v\n"+
		"noL2Recalcs: %+v\n"+
		"totalDepositedAmount: %d\n"+
		"totalWithdrawnAmount: %d\n"+
		"rollupWithMoreRecentProof: %d\n"+
		"nrTransferTransactions: %d\n",
		o.simulation.l1Network.Stats.nrMiners,
		o.l1Height,
		o.l2Height,
		o.simulation.l1Network.Stats.totalL1Blocks,
		o.simulation.l1Network.Stats.totalL2Blocks,
		o.l2RollupCountInHeaders,
		o.l2RollupCountInL1Blocks,
		o.l2RollupTxCountInL1Blocks,
		o.simulation.l1Network.Stats.maxRollupsPerBlock,
		o.simulation.l1Network.Stats.nrEmptyBlocks,
		o.simulation.l1Network.Stats.totalL2Blocks,
		o.simulation.l1Network.Stats.noL1Reorgs,
		o.simulation.l1Network.Stats.noL2Recalcs,
		o.simulation.l1Network.Stats.totalDepositedAmount,
		o.simulation.l1Network.Stats.totalWithdrawnAmount,
		o.simulation.l1Network.Stats.rollupWithMoreRecentProof,
		o.simulation.l1Network.Stats.nrTransferTransactions,
	)
}
