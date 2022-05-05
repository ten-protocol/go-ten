package simulation

import (
	"fmt"
	"sync"
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
		t.Logf("Node Heights: %v", heights)
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

	// process the blockchain of each node in parallel to minimize the difference between them since they are still running
	heights := make([]uint64, len(s.ObscuroNodes))
	var wg sync.WaitGroup
	for i := range s.ObscuroNodes {
		wg.Add(1)
		go checkBlockchainOfObscuroNode(t, s.ObscuroNodes[i], minHeight, maxL1Height, s, &wg, heights, i)
	}
	wg.Wait()
	min, max := minMax(heights)
	// This checks that all the nodes are in sync. When a node falls behind with processing blocks it might highlight a problem.
	if max-min > max/10 {
		t.Errorf("There is a problem with the Obscuro chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}
}

func checkBlockchainOfEthereumNode(t *testing.T, node ethclient.EthClient, minHeight uint64, s *Simulation) uint64 {
	head := node.FetchHeadBlock()
	height := head.NumberU64()

	if height < minHeight {
		t.Errorf("Node %d. There were only %d blocks mined. Expected at least: %d.", obscurocommon.ShortAddress(node.Info().ID), height, minHeight)
	}

	deposits, rollups, totalDeposited, blockCount := extractDataFromEthereumChain(head, node, s)
	s.Stats.TotalL1Blocks = uint64(blockCount)

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
		t.Errorf("Node %d. Efficiency in L1 is %f. Expected:%f. Number: %d.", obscurocommon.ShortAddress(node.Info().ID), efficiency, s.Params.L1EfficiencyThreshold, height)
	}

	// compare the number of reorgs for this node against the height
	reorgs := s.Stats.NoL1Reorgs[node.Info().ID]
	eff := float64(reorgs) / float64(height)
	if eff > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d. The number of reorgs is too high: %d. ", obscurocommon.ShortAddress(node.Info().ID), reorgs)
	}
	return height
}

func extractDataFromEthereumChain(head *types.Block, node ethclient.EthClient, s *Simulation) ([]common.Hash, []obscurocommon.L2RootHash, uint64, int) {
	deposits := make([]common.Hash, 0)
	rollups := make([]obscurocommon.L2RootHash, 0)
	totalDeposited := uint64(0)

	blockchain := node.BlocksBetween(obscurocommon.GenesisBlock, head)
	for _, block := range blockchain {
		for _, tx := range block.Transactions() {
			t := s.Params.TxHandler.UnPackTx(tx)
			if t == nil {
				continue
			}
			switch t.TxType {
			case obscurocommon.DepositTx:
				deposits = append(deposits, tx.Hash())
				totalDeposited += t.Amount
			case obscurocommon.RollupTx:
				r := nodecommon.DecodeRollupOrPanic(t.Rollup)
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
	return deposits, rollups, totalDeposited, len(blockchain)
}

// MAX_BLOCK_DELAY the maximum an Obscuro node can fall behind
const MAX_BLOCK_DELAY = 5 // nolint:revive,stylecheck

func checkBlockchainOfObscuroNode(t *testing.T, node *host.Node, minObscuroHeight uint64, maxEthereumHeight uint64, s *Simulation, wg *sync.WaitGroup, heights []uint64, i int) uint64 {
	l1Height := uint64(node.DB().GetCurrentBlockHead().Number.Int64())

	// check that the L1 view is consistent with the L1 network.
	// We cast to int64 to avoid an overflow when l1Height is greater than maxEthereumHeight (due to additional blocks
	// produced since maxEthereumHeight was calculated from querying all L1 nodes - the simulation is still running, so
	// new blocks might have been added in the meantime).
	if int64(maxEthereumHeight)-int64(l1Height) > MAX_BLOCK_DELAY {
		t.Errorf("Obscuro node %d fell behind %d blocks.", obscurocommon.ShortAddress(node.ID), maxEthereumHeight-l1Height)
	}

	// check that the height of the Rollup chain is higher than a minimum expected value.
	h := node.DB().GetCurrentRollupHead()
	if h == nil {
		panic(fmt.Sprintf("Node %d has no head rollup recorded.\n", obscurocommon.ShortAddress(node.ID)))
	}
	l2Height := h.Number
	if l2Height < minObscuroHeight {
		t.Errorf("There were only %d rollups mined on node %d. Expected at least: %d.", l2Height, obscurocommon.ShortAddress(node.ID), minObscuroHeight)
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
	transfers, withdrawals := s.TxInjector.GetL2Transactions()
	notFoundTransfers := 0
	for _, tx := range transfers {
		if l2tx := node.Enclave.GetTransaction(tx.Hash()); l2tx == nil {
			notFoundTransfers++
		}
	}
	if notFoundTransfers > 0 {
		t.Errorf("Node %d - %d out of %d Transfer Txs not found in the enclave", obscurocommon.ShortAddress(node.ID), notFoundTransfers, len(transfers))
	}

	notFoundWithdrawals := 0
	for _, tx := range withdrawals {
		if l2tx := node.Enclave.GetTransaction(tx.Hash()); l2tx == nil {
			notFoundWithdrawals++
		}
	}
	if notFoundWithdrawals > 0 {
		t.Errorf("Node %d - %d out of %d Withdrawal Txs not found in the enclave", obscurocommon.ShortAddress(node.ID), notFoundWithdrawals, len(withdrawals))
	}

	totalSuccessfullyWithdrawn, numberOfWithdrawalRequests := extractWithdrawals(node)

	// sanity check number of withdrawal transaction
	if numberOfWithdrawalRequests > len(s.TxInjector.GetL2WithdrawalRequests()) {
		t.Errorf("found more transactions in the blockchain than the generated by the tx manager")
	}

	injectorDepositedAmt := uint64(0)
	for _, tx := range s.TxInjector.GetL1Transactions() {
		injectorDepositedAmt += tx.Amount
	}

	// expected condition : some Txs (stats) did not make it to the blockchain
	// best condition : all Txs (stats) were issue and consumed in the blockchain
	// can't happen : sum of headers withdraws greater than issued Txs (stats)
	if totalSuccessfullyWithdrawn > s.Stats.TotalWithdrawalRequestedAmount {
		t.Errorf("The amount withdrawn %d is exceeds the actual amount requested %d", totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// sanity check that the injected withdrawals were mostly executed
	if totalSuccessfullyWithdrawn < s.Stats.TotalWithdrawalRequestedAmount/2 {
		t.Errorf("The amount withdrawn %d is far smaller than the amount requested %d. Something is probably wrong.", totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
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

	heights[i] = l2Height
	wg.Done()
	return l2Height
}

func extractWithdrawals(node *host.Node) (totalSuccessfullyWithdrawn uint64, numberOfWithdrawalRequests int) {
	head := node.DB().GetCurrentRollupHead()
	if head == nil {
		panic("the current head should not be nil")
	}
	// sum all the withdrawals by traversing the node headers from Head to Genesis
	for r := head; ; r = node.DB().GetRollupHeader(r.ParentHash) {
		if r != nil && r.Number == obscurocommon.L1GenesisHeight {
			return
		}
		if r == nil {
			panic(fmt.Sprintf("Reached a missing rollup on node %d", obscurocommon.ShortAddress(node.ID)))
		}
		for _, w := range r.Withdrawals {
			totalSuccessfullyWithdrawn += w.Amount
			numberOfWithdrawalRequests++
		}
	}
}
