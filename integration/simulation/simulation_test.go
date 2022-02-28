package simulation

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	"golang.org/x/sync/errgroup"

	obscuroCommon "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

func TestSimulation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	// create a folder specific for the test
	err := os.MkdirAll("../.build/simulations/", 0o700)
	if err != nil {
		panic(err)
	}
	f, err := os.CreateTemp("../.build/simulations", "simulation-result-*.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	common.SetLog(f)

	txGenerator := NewTransactionGenerator(5)

	numberOfNodes := 10
	simulationTime := 15
	avgBlockDuration := uint64(20_000)
	avgLatency := avgBlockDuration / 15
	avgGossipPeriod := avgBlockDuration / 3

	stats := NewStats(numberOfNodes, simulationTime, avgBlockDuration, avgLatency, avgGossipPeriod)

	l1Network, l2Network := RunSimulation(txGenerator, numberOfNodes, simulationTime, avgBlockDuration, avgBlockDuration/15, avgBlockDuration/3, stats)

	checkBlockchainValidity(t, l1Network, l2Network, txGenerator)

	fmt.Println("Simulation ended...")
	fmt.Printf("%+v\n", stats)
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, txManager *TransactionManager) {
	// TODO check all nodes are the same height ?
	// pick one node to draw height
	enclaveNode := l2Network.nodes[0].Enclave
	l1Node := l1Network.nodes[0]
	l1Height := enclaveNode.L1Height()
	l1HeightHash := enclaveNode.L1HeightHash()
	l2Height := enclaveNode.L2Height()

	validateL1(t, l1Network.Stats, l1Height, l1HeightHash, l1Node)
	totalWithdrawn := validateL2Stats(t, l1Network.Stats, l2Height, txManager)
	validateL2StateStats(t, l2Network, l1Network.Stats, totalWithdrawn, txManager.wallets)
	validateL2State(t, l2Network.nodes, txManager)
}

// validateL2State tests that all transaction in the transaction Manager are found in the blockchain state of each node
func validateL2State(t *testing.T, nodes []*obscuro_node.Node, txManager *TransactionManager) {
	// Parallelize this check
	var nGroup errgroup.Group

	// Check the balance of all nodes adds up with the balance in the stats
	for _, node := range nodes {
		closureNode := node
		nGroup.Go(func() error {
			// all transactions should exist on every node
			for _, transaction := range txManager.GetL2Transactions() {
				tx, found := closureNode.Enclave.GetTransaction(transaction.ID)
				if !found {
					return fmt.Errorf("unable to find transaction: %v", tx) // nolint:goerr113
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
const L1EfficiencyThreashold = 0.2
const L2EfficiencyThreashold = 0.3

// validateL1 does a sanity check on the mock implementation of the L1
func validateL1(t *testing.T, stats *Stats, l1Height int, l1HeightHash common.L1RootHash, node *ethereum_mock.Node) {
	deposits := make([]uuid.UUID, 0)
	rollups := make([]common.L2RootHash, 0)
	stats.l1Height = l1Height
	totalDeposited := uint64(0)

	l1Block, found := node.Resolver.Resolve(l1HeightHash)
	if !found {
		t.Errorf("expected l1 height block not found")
	}

	blockchain := ethereum_mock.BlocksBetween(&common.GenesisBlock, l1Block, node.Resolver)
	for _, block := range blockchain {
		for _, tx := range block.Transactions {
			switch tx.TxType {
			case common.DepositTx:
				deposits = append(deposits, tx.ID)
				totalDeposited += tx.Amount
			case common.RollupTx:
				r := obscuroCommon.DecodeRollup(tx.Rollup)
				rollups = append(rollups, r.Hash())
				if common.IsBlockAncestor(r.Header.L1Proof, block, node.Resolver) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					stats.NewRollup(r)
				}
			default:
				panic("unknown transaction type")
			}
		}
	}

	if len(common.FindDups(deposits)) > 0 {
		dups := common.FindDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(common.FindRollupDups(rollups)) > 0 {
		dups := common.FindRollupDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != stats.totalDepositedAmount {
		t.Errorf("Deposit amounts don't match. Found %d , expected %d", totalDeposited, stats.totalDepositedAmount)
	}

	efficiency := float64(stats.totalL1Blocks-stats.l1Height) / float64(stats.totalL1Blocks)
	if efficiency > L1EfficiencyThreashold {
		t.Errorf("Efficiency in L1 is %f. Expected:%f", efficiency, L1EfficiencyThreashold)
	}

	// todo
	for nodeID, reorgs := range stats.noL1Reorgs {
		eff := float64(reorgs) / float64(stats.l1Height)
		if eff > L1EfficiencyThreashold {
			t.Errorf("Efficiency for node %d in L1 is %f. Expected:%f", nodeID, eff, L1EfficiencyThreashold)
		}
	}
}

func validateL2Stats(t *testing.T, stats *Stats, l2Height int, txManager *TransactionManager) uint64 {
	// get non-failed transactions
	transactions := txManager.GetL2Transactions()
	withdrawalRequests := txManager.GetL2WithdrawalRequests()

	// todo - check that proofs are on the canonical chain

	// Sum all withdraw transactions AND
	txSumWithdrawalsRequested := sumWithdrawals(withdrawalRequests)

	// are there duplicated transaction Ids ?
	txIds := make([]uuid.UUID, len(transactions))
	for i, transaction := range transactions {
		txIds[i] = transaction.ID
	}
	if dups := common.FindDups(txIds); len(dups) > 0 {
		t.Errorf("Found L2 txs duplicates: %v", dups)
	}

	// are the number of txs sent by the generator the same as the captured by the enclave stats ?
	if len(transactions) != stats.nrTransferTransactions {
		t.Errorf("Nr of transfers don't match. Found %d , expected %d", len(transactions), stats.nrTransferTransactions)
	}

	if txSumWithdrawalsRequested != stats.totalWithdrawnAmount {
		t.Errorf("The amount withdrawn %d exceeds the actual amount requested %d", txSumWithdrawalsRequested, stats.totalWithdrawnAmount)
	}

	efficiency := float64(stats.totalL2Blocks-l2Height) / float64(stats.totalL2Blocks)
	if efficiency > L2EfficiencyThreashold {
		t.Errorf("Efficiency in L2 is %f. Expected:%f", efficiency, L2EfficiencyThreashold)
	}

	return txSumWithdrawalsRequested
}

func sumWithdrawals(w []obscuroCommon.Withdrawal) uint64 {
	sum := uint64(0)
	for _, r := range w {
		sum += r.Amount
	}
	return sum
}

func validateL2StateStats(t *testing.T, l2Network L2NetworkCfg, s *Stats, totalWithdrawn uint64, wallets []wallet_mock.Wallet) {
	finalAmount := s.totalDepositedAmount - totalWithdrawn

	// Parallelize this check
	var nGroup errgroup.Group

	// Check the balance of all nodes adds up with the balance in the stats
	for _, node := range l2Network.nodes {
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
