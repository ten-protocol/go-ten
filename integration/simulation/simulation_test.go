package simulation

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"

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

	logFile := setupTestLog("../.build/simulations/")
	defer logFile.Close()

	// define core test parameters
	numberOfNodes := 10
	simulationTimeSecs := 15                // in seconds
	avgBlockDurationUSecs := uint64(20_000) // in u seconds 1 sec = 1e6 usecs
	avgLatency := avgBlockDurationUSecs / 15
	avgGossipPeriod := avgBlockDurationUSecs / 3

	// define network params
	stats := NewStats(numberOfNodes)
	l1NetworkConfig := NewL1Network(avgLatency, stats)
	l2NetworkCfg := NewL2Network(avgLatency)

	// define instances of the simulation mechanisms
	txManager := NewTransactionManager(5, l1NetworkConfig, l2NetworkCfg, avgBlockDurationUSecs, stats)
	simulation := NewSimulation(
		numberOfNodes,
		l1NetworkConfig,
		l2NetworkCfg,
		avgBlockDurationUSecs,
		avgGossipPeriod,
		stats,
	)

	// execute the simulation
	simulation.Start(txManager, simulationTimeSecs)

	// run tests
	checkBlockchainValidity(t, txManager, simulation)

	t.Logf("%+v\n", stats)
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func checkBlockchainValidity(t *testing.T, txManager *TransactionManager, network *Simulation) {
	// TODO check all nodes are the same height ?
	// pick one node to draw height
	l1Node := network.l1Network.nodes[0]
	obscuroNode := network.l2Network.nodes[0]
	currentBlockHead := obscuroNode.Headers().GetCurrentBlockHead()
	currentRollupHead := obscuroNode.Headers().GetCurrentRollupHead()

	l1Height := currentBlockHead.Height
	l1HeightHash := currentBlockHead.ID
	l2Height := currentRollupHead.Height

	// ensure the L1 blocks are valid
	validateL1(t, network.l1Network.Stats, l1Height, &l1HeightHash, l1Node)

	// ensure the validity of l1 vs l2 stats
	validateL1L2Stats(t, obscuroNode, network.l1Network.Stats)

	// ensure the generated withdrawal stats match the l2 blockchain state (withdrawals)
	totalWithdrawn := validateL2WithdrawalStats(t, obscuroNode, network.l1Network.Stats, l2Height, txManager)

	// ensure that each node has the expected total balance computed above
	validateL2NodeBalances(t, network.l2Network, network.l1Network.Stats, totalWithdrawn, txManager.wallets)

	// ensure that each node can fetch each of the generated transactions
	validateL2TxsExist(t, network.l2Network.nodes, txManager)
}

// validateL1L2Stats validates blockchain wide properties between L1 and the L2
func validateL1L2Stats(t *testing.T, node *obscuro_node.Node, stats *Stats) {
	l1HeaderCount := uint(0)
	for header := node.Headers().GetCurrentBlockHead(); header != nil; header = node.Headers().GetBlockHeader(header.Parent) {
		l1HeaderCount++
	}
	l2HeaderCount := uint(0)
	for header := node.Headers().GetCurrentRollupHead(); header != nil; header = node.Headers().GetRollupHeader(header.Parent) {
		l2HeaderCount++
	}

	if l1HeaderCount > stats.totalL1Blocks || l2HeaderCount > stats.totalL2Blocks {
		t.Errorf("should not have more blocks/rollups in stats than in the node header "+
			"- Blocks: Header %d, Stats %d - Rollups: Header %d, Stats %d ",
			l1HeaderCount,
			stats.totalL1Blocks,
			l2HeaderCount,
			stats.totalL2Blocks,
		)
	}

	efficiency := float64(l1HeaderCount-l2HeaderCount) / float64(l1HeaderCount)
	if efficiency > L2ToL1EfficiencyThreshold {
		t.Errorf("L2 to L1 Efficiency is %f. Expected:%f", efficiency, L2ToL1EfficiencyThreshold)
	}

	t.Logf("There was %d L1 blocks and %d L2 Blocks\n", stats.totalL1Blocks, stats.totalL2Blocks)
	t.Logf("Node %d Header had %d l2 blocks in the header\n", l1HeaderCount, l2HeaderCount)
}

// validateL2TxsExist tests that all transaction in the transaction Manager are found in the blockchain state of each node
func validateL2TxsExist(t *testing.T, nodes []*obscuro_node.Node, txManager *TransactionManager) {
	// Parallelize this check
	var nGroup errgroup.Group

	// Create a go routine to check each node
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
const (
	L1EfficiencyThreshold     = 0.2
	L2EfficiencyThreshold     = 0.3
	L2ToL1EfficiencyThreshold = 0.3
)

// validateL1 does a sanity check on the mock implementation of the L1
func validateL1(t *testing.T, stats *Stats, l1Height uint, l1HeightHash *common.L1RootHash, node *ethereum_mock.Node) {
	deposits := make([]uuid.UUID, 0)
	rollups := make([]common.L2RootHash, 0)
	totalDeposited := uint64(0)

	l1Block, found := node.Resolver.Resolve(*l1HeightHash)
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
			case common.RequestSecretTx:
			case common.StoreSecretTx:
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
func validateL2WithdrawalStats(t *testing.T, node *obscuro_node.Node, stats *Stats, l2Height uint, txManager *TransactionManager) uint64 {
	headerWithdrawalSum := uint64(0)
	headerWithdrawalTxCount := 0

	// todo - check that proofs are on the canonical chain
	// sum all the withdrawals by traversing the node headers from Head to Genesis
	for header := node.Headers().GetCurrentRollupHead(); header != nil; header = node.Headers().GetRollupHeader(header.Parent) {
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

	// TODO - there should be an efficiency test between blocks and rollups
	// you should not have % difference between the # of rollups and the # of blocks

	efficiency := float64(stats.totalL2Blocks-l2Height) / float64(stats.totalL2Blocks)
	if efficiency > L2EfficiencyThreshold {
		t.Errorf("Efficiency in L2 is %f. Expected:%f", efficiency, L2EfficiencyThreshold)
	}

	return headerWithdrawalSum
}

func validateL2NodeBalances(t *testing.T, l2Network *L2NetworkCfg, s *Stats, totalWithdrawn uint64, wallets []wallet_mock.Wallet) {
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

func setupTestLog(baseDir string) *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(baseDir, 0o700)
	if err != nil {
		panic(err)
	}
	f, err := os.CreateTemp(baseDir, "simulation-result-*.txt")
	if err != nil {
		panic(err)
	}
	log.SetLog(f)
	return f
}
