package simulation

import (
	"testing"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// After a simulation has run, check as much as possible that the outputs of the simulation are expected.
// For example, all injected transactions were processed correctly, the height of the rollup chain is a function of the total
// time of the simulation and the average block duration, that all Obscuro nodes are roughly in sync, etc
func checkNetworkValidity(t *testing.T, s *Simulation) {
	// ensure L1 and L2 txs were issued
	if len(s.TxInjector.l1Transactions) == 0 || len(s.TxInjector.l2Transactions) == 0 {
		t.Error("Not enough transactions issued")
	}

	l1MaxHeight := checkEthereumBlockchainValidity(t, s)
	checkObscuroBlockchainValidity(t, s, l1MaxHeight)
}

// checkEthereumBlockchainValidity: sanity check on the mock implementation of the L1 on all nodes
// - minimum height - the chain has a minimum number of blocks
// - check height is similar across all Mock ethereum nodes
// - check no duplicate txs
// - check efficiency - no of created blocks/ height
// - noReorgs
func checkEthereumBlockchainValidity(t *testing.T, s *Simulation) uint64 {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	heights := make([]uint64, len(s.EthClients))
	for i, node := range s.EthClients {
		heights[i] = checkBlockchainOfEthereumNode(t, node, minHeight, s)
	}

	min, max := minMax(heights)
	if max-min > max/10 {
		t.Errorf("There is a problem with the Mock ethereum chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}

	return max
}

// checkObscuroBlockchainValidity - perform the following checks
// - minimum height - the chain has a minimum number of rollups
// - check height is similar
// - check no duplicate txs
// - check efficiency - no of created blocks/ height
// - check amount in the system
// - check withdrawals/deposits
func checkObscuroBlockchainValidity(t *testing.T, s *Simulation, maxL1Height uint64) {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	heights := make([]uint64, len(s.ObscuroNodes))
	for i, node := range s.ObscuroNodes {
		heights[i] = checkBlockchainOfObscuroNode(t, node, minHeight, maxL1Height, s)
	}

	min, max := minMax(heights)
	// This checks that all the nodes are in sync. When a node falls behind with processing blocks it might highlight a problem.
	if max-min > max/10 {
		t.Errorf("There is a problem with the Obscuro chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}
}

func checkBlockchainOfEthereumNode(t *testing.T, node ethclient.Client, minHeight uint64, s *Simulation) uint64 {
	head, height := node.FetchHeadBlock()

	if height < minHeight {
		t.Errorf("Node %d. There were only %d blocks mined. Expected at least: %d.", obscurocommon.ShortAddress(node.Info().ID), height, minHeight)
	}

	deposits, rollups, totalDeposited := extractDataFromEthereumChain(head, node, s)

	if len(obscurocommon.FindHashDups(deposits)) > 0 {
		dups := obscurocommon.FindHashDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(obscurocommon.FindRollupDups(rollups)) > 0 {
		dups := obscurocommon.FindRollupDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != s.Stats.TotalDepositedAmount {
		t.Errorf("Node %d. Deposit amounts don't match. Found %d , expected %d", obscurocommon.ShortAddress(node.Info().ID), totalDeposited, s.Stats.TotalDepositedAmount)
	}

	efficiency := float64(s.Stats.TotalL1Blocks-height) / float64(s.Stats.TotalL1Blocks)
	if efficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d. Efficiency in L1 is %f. Expected:%f. Height: %d.", obscurocommon.ShortAddress(node.Info().ID), efficiency, s.Params.L1EfficiencyThreshold, height)
	}

	// compare the number of reorgs for this node against the height
	reorgs := s.Stats.NoL1Reorgs[node.Info().ID]
	eff := float64(reorgs) / float64(height)
	if eff > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d. The number of reorgs is too high: %d. ", obscurocommon.ShortAddress(node.Info().ID), reorgs)
	}
	return height
}

func extractDataFromEthereumChain(head *types.Block, node ethclient.Client, s *Simulation) ([]common.Hash, []obscurocommon.L2RootHash, uint64) {
	deposits := make([]common.Hash, 0)
	rollups := make([]obscurocommon.L2RootHash, 0)
	totalDeposited := uint64(0)

	blockchain := node.BlocksBetween(obscurocommon.GenesisBlock, head)
	for _, block := range blockchain {
		for _, tr := range block.Transactions() {
			tx := obscurocommon.TxData(tr)
			switch tx.TxType {
			case obscurocommon.DepositTx:
				deposits = append(deposits, tr.Hash())
				totalDeposited += tx.Amount
			case obscurocommon.RollupTx:
				r := nodecommon.DecodeRollupOrPanic(tx.Rollup)
				rollups = append(rollups, r.Hash())
				if node.IsBlockAncestor(block, r.Header.L1Proof) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					s.Stats.NewRollup(node.Info().ID, r)
				}
			case obscurocommon.RequestSecretTx:
			case obscurocommon.StoreSecretTx:
			}
		}
	}
	return deposits, rollups, totalDeposited
}

// MAX_BLOCK_DELAY the maximum an Obscuro node can fall behind
const MAX_BLOCK_DELAY = 5 // nolint:revive,stylecheck

func checkBlockchainOfObscuroNode(t *testing.T, node *host.Node, minObscuroHeight uint64, maxEthereumHeight uint64, s *Simulation) uint64 {
	l1Height := node.DB().GetCurrentBlockHead().Height

	// check that the L1 view is consistent with the L1 network.
	// We cast to int64 to avoid an overflow when l1Height is greater than maxEthereumHeight (due to additional blocks
	// produced since maxEthereumHeight was calculated from querying all L1 nodes - the simulation is still running, so
	// new blocks might have been added in the meantime).
	if int64(maxEthereumHeight)-int64(l1Height) > MAX_BLOCK_DELAY {
		t.Errorf("Obscuro node %d fell behind %d blocks.", obscurocommon.ShortAddress(node.ID), maxEthereumHeight-l1Height)
	}

	// check that the height of the Rollup chain is higher than a minimum expected value.
	l2Height := node.DB().GetCurrentRollupHead().Height
	if l2Height < minObscuroHeight {
		t.Errorf("There were only %d blocks mined on node %d. Expected at least: %d.", l2Height, obscurocommon.ShortAddress(node.ID), minObscuroHeight)
	}

	totalL2Blocks := s.Stats.NoL2Blocks[node.ID]
	// in case the blockchain has advanced above what was collected, there is no longer a point to this check
	if l2Height <= totalL2Blocks {
		efficiencyL2 := float64(totalL2Blocks-l2Height) / float64(totalL2Blocks)
		if efficiencyL2 > s.Params.L2EfficiencyThreshold {
			t.Errorf("Node %d. Efficiency in L2 is %f. Expected:%f", obscurocommon.ShortAddress(node.ID), efficiencyL2, s.Params.L2EfficiencyThreshold)
		}
	}

	// check that the pobi protocol doesn't waste too many blocks.
	// todo- find the block where the genesis was published)
	efficiency := float64(l1Height-l2Height) / float64(l1Height)
	if efficiency > s.Params.L2ToL1EfficiencyThreshold {
		t.Errorf("L2 to L1 Efficiency is %f. Expected:%f", efficiency, s.Params.L2ToL1EfficiencyThreshold)
	}

	// check that all expected transactions were included.
	for _, transaction := range s.TxInjector.GetL2Transactions() {
		l2tx := node.Enclave.GetTransaction(transaction.Hash())
		if l2tx == nil {
			t.Errorf("node %d, unable to find transaction: %+v", obscurocommon.ShortAddress(node.ID), transaction)
		}
	}

	totalSuccessfullyWithdrawn, numberOfWithdrawalRequests := extractWithdrawals(node)

	// sanity check number of withdrawal transaction
	if numberOfWithdrawalRequests > len(s.TxInjector.GetL2WithdrawalRequests()) {
		t.Errorf("found more transactions in the blockchain than the generated by the tx manager")
	}

	// expected condition : some Txs (stats) did not make it to the blockchain
	// best condition : all Txs (stats) were issue and consumed in the blockchain
	// can't happen : sum of headers withdraws greater than issued Txs (stats)
	if totalSuccessfullyWithdrawn > s.Stats.TotalWithdrawalRequestedAmount {
		t.Errorf("The amount withdrawn %d is not the same as the actual amount requested %d", totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// check that the sum of all balances matches the total amount of money that must be in the system
	totalAmountInSystem := s.Stats.TotalDepositedAmount - totalSuccessfullyWithdrawn
	total := uint64(0)
	for _, wallet := range s.TxInjector.wallets {
		total += node.Enclave.Balance(wallet.Address)
	}
	if total != totalAmountInSystem {
		t.Errorf("the amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", obscurocommon.ShortAddress(node.ID), total, totalAmountInSystem)
	}
	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// (execute deposits and transactions and compare to the state in the rollup)

	return l2Height
}

func extractWithdrawals(node *host.Node) (totalSuccessfullyWithdrawn uint64, numberOfWithdrawalRequests int) {
	// sum all the withdrawals by traversing the node headers from Head to Genesis
	for r := node.DB().GetCurrentRollupHead(); r != nil; r = node.DB().GetRollupHeader(r.Parent) {
		for _, w := range r.Withdrawals {
			totalSuccessfullyWithdrawn += w.Amount
			numberOfWithdrawalRequests++
		}
	}
	return
}
