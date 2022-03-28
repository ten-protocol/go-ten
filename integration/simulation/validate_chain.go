package simulation

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
	"golang.org/x/sync/errgroup"
)

func checkBlockchainValidity(t *testing.T, s *Simulation) {
	// TODO check all nodes are the same height ?
	// pick one node to draw height
	l1Node := s.MockEthNodes[0]
	obscuroNode := s.InMemObscuroNodes[0]
	currentBlockHead := obscuroNode.DB().GetCurrentBlockHead()
	currentRollupHead := obscuroNode.DB().GetCurrentRollupHead()

	l1Height := currentBlockHead.Height
	l1HeightHash := currentBlockHead.ID
	l2Height := currentRollupHead.Height

	// ensure the L1 blocks are valid
	validateL1(t, s.Stats, l1Height, &l1HeightHash, l1Node)

	// ensure the validity of l1 vs l2 stats
	validateL1L2Stats(t, obscuroNode, s.Stats)

	// ensure the generated withdrawal stats match the l2 blockchain state (withdrawals)
	totalWithdrawn := validateL2WithdrawalStats(t, obscuroNode, s.Stats, l2Height, s.TxInjector)

	// ensure that each node has the expected total balance computed above
	validateL2NodeBalances(t, s.InMemObscuroNodes, s.Stats, totalWithdrawn, s.TxInjector.wallets)

	// ensure that each node can fetch each of the generated transactions
	validateL2TxsExist(t, s.InMemObscuroNodes, s.TxInjector)
}

// validateL1L2Stats validates blockchain wide properties between L1 and the L2
func validateL1L2Stats(t *testing.T, node *host.Node, stats *Stats) {
	l1Height := obscurocommon.L1GenesisHeight
	for header := node.DB().GetCurrentBlockHead(); header != nil && header.ID != obscurocommon.GenesisHash; header = node.DB().GetBlockHeader(header.Parent) {
		l1Height++
	}
	l2Height := obscurocommon.L2GenesisHeight
	for header := node.DB().GetCurrentRollupHead(); header != nil && header.ID != obscurocommon.GenesisHash; header = node.DB().GetRollupHeader(header.Parent) {
		l2Height++
	}

	// todo - figure out why +1
	if l1Height != node.DB().GetCurrentBlockHead().Height+1 {
		t.Errorf("unexpected block height. expected %d, got %d", l1Height, node.DB().GetCurrentBlockHead().Height)
	}

	// todo - figure out why +1
	if l2Height != node.DB().GetCurrentRollupHead().Height+1 {
		t.Errorf("unexpected rollup height. expected %d, got %d", l2Height, node.DB().GetCurrentRollupHead().Height)
	}

	if l1Height > stats.totalL1Blocks || l2Height > stats.totalL2Blocks {
		t.Errorf("should not have more blocks/rollups in stats than in the node header "+
			"- Blocks: Header %d, Stats %d - Rollups: Header %d, Stats %d ",
			l1Height,
			stats.totalL1Blocks,
			l2Height,
			stats.totalL2Blocks,
		)
	}

	efficiency := float64(l1Height-l2Height) / float64(l1Height)
	if efficiency > L2ToL1EfficiencyThreshold {
		t.Errorf("L2 to L1 Efficiency is %f. Expected:%f", efficiency, L2ToL1EfficiencyThreshold)
	}
}

// validateL2TxsExist tests that all transaction in the transaction Manager are found in the blockchain state of each node
func validateL2TxsExist(t *testing.T, nodes []*host.Node, txManager *TransactionInjector) {
	// Parallelize this check
	var nGroup errgroup.Group

	// Create a go routine to check each node
	for _, node := range nodes {
		closureNode := node
		nGroup.Go(func() error {
			// all transactions should exist on every node
			for _, transaction := range txManager.GetL2Transactions() {
				l2tx := closureNode.Enclave.GetTransaction(transaction.Hash())
				if l2tx == nil {
					return fmt.Errorf("node %d, unable to find transaction: %+v", closureNode.ID, transaction) // nolint:goerr113
				}
			}
			return nil
		})
	}
	if err := nGroup.Wait(); err != nil {
		t.Error(err)
	}
}

// For this simulation, this represents an acceptable "dead blocks" percentage.
// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
// We test the results against this threshold to catch eventual protocol errors.
const (
	L1EfficiencyThreshold     = 0.2
	L2EfficiencyThreshold     = 0.3
	L2ToL1EfficiencyThreshold = 0.32
)

// validateL1 does a sanity check on the mock implementation of the L1
func validateL1(t *testing.T, stats *Stats, l1Height uint64, l1HeightHash *obscurocommon.L1RootHash, node *ethereum_mock.Node) {
	deposits := make([]common.Hash, 0)
	rollups := make([]obscurocommon.L2RootHash, 0)
	totalDeposited := uint64(0)

	l1Block, found := node.Resolver.FetchBlock(*l1HeightHash)
	if !found {
		t.Errorf("expected l1 height block not found")
	}

	blockchain := ethereum_mock.BlocksBetween(obscurocommon.GenesisBlock, l1Block, node.Resolver)
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
				if node.Resolver.IsBlockAncestor(block, r.Header.L1Proof) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					stats.NewRollup(r)
				}
			case obscurocommon.RequestSecretTx:
			case obscurocommon.StoreSecretTx:
			}
		}
	}

	if len(obscurocommon.FindHashDups(deposits)) > 0 {
		dups := obscurocommon.FindHashDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(obscurocommon.FindRollupDups(rollups)) > 0 {
		dups := obscurocommon.FindRollupDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != stats.totalDepositedAmount {
		t.Errorf("Deposit amounts don't match. Found %d , expected %d", totalDeposited, stats.totalDepositedAmount)
	}

	efficiency := float64(stats.totalL1Blocks-l1Height) / float64(stats.totalL1Blocks)
	if efficiency > L1EfficiencyThreshold {
		t.Errorf("Efficiency in L1 is %f. Expected:%f", efficiency, L1EfficiencyThreshold)
	}

	// todo
	for nodeID, reorgs := range stats.noL1Reorgs {
		eff := float64(reorgs) / float64(l1Height)
		if eff > L1EfficiencyThreshold {
			t.Errorf("Efficiency for node %d in L1 is %f. Expected:%f", nodeID, eff, L1EfficiencyThreshold)
		}
	}
}

// validateL2WithdrawalStats checks the withdrawal requests by
// comparing the stats of the generated transactions with the withdrawals on the node headers
func validateL2WithdrawalStats(t *testing.T, node *host.Node, stats *Stats, l2Height uint64, txManager *TransactionInjector) uint64 {
	headerWithdrawalSum := uint64(0)
	headerWithdrawalTxCount := 0

	// todo - check that proofs are on the canonical chain
	// sum all the withdrawals by traversing the node headers from Head to Genesis
	for header := node.DB().GetCurrentRollupHead(); header != nil; header = node.DB().GetRollupHeader(header.Parent) {
		for _, w := range header.Withdrawals {
			headerWithdrawalSum += w.Amount
			headerWithdrawalTxCount++
		}
	}

	// get all generated withdrawal txs
	if headerWithdrawalTxCount > len(txManager.GetL2WithdrawalRequests()) {
		t.Errorf("found more transactions in the blockchain than the generated by the tx manager")
	}

	// expected condition : some Txs (stats) did not make it to the blockchain
	// best condition : all Txs (stats) were issue and consumed in the blockchain
	// can't happen : sum of headers withdraws greater than issued Txs (stats)
	if headerWithdrawalSum > stats.totalWithdrawnAmount {
		t.Errorf("The amount withdrawn %d is not the same as the actual amount requested %d", headerWithdrawalSum, stats.totalWithdrawnAmount)
	}

	// you should not have % difference between the # of rollups and the # of blocks
	efficiency := float64(stats.totalL2Blocks-l2Height) / float64(stats.totalL2Blocks)
	if efficiency > L2EfficiencyThreshold {
		t.Errorf("Efficiency in L2 is %f. Expected:%f", efficiency, L2EfficiencyThreshold)
	}

	return headerWithdrawalSum
}

func validateL2NodeBalances(t *testing.T, l2Nodes []*host.Node, s *Stats, totalWithdrawn uint64, wallets []wallet_mock.Wallet) {
	finalAmount := s.totalDepositedAmount - totalWithdrawn

	// Parallelize this check
	var nGroup errgroup.Group

	// Check the balance of all nodes adds up with the balance in the stats
	for _, node := range l2Nodes {
		closureNode := node
		nGroup.Go(func() error {
			// add up all balances
			total := uint64(0)
			for _, wallet := range wallets {
				total += closureNode.Enclave.Balance(wallet.Address)
			}
			if total != finalAmount {
				return fmt.Errorf("the amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", closureNode.ID, total, finalAmount) // nolint:goerr113
			}
			return nil
		})
	}
	if err := nGroup.Wait(); err != nil {
		t.Error(err)
	}
	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// walk the blocks in reverse direction, execute deposits and transactions and compare to the state in the rollup
}
