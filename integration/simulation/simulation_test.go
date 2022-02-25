package simulation

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
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

	blockDuration := uint64(20_000)
	l1netw, l2netw, wallets := RunSimulation(5, 10, 15, blockDuration, blockDuration/15, blockDuration/3)
	checkBlockchainValidity(t, l1netw, l2netw, wallets)
	stats := l1netw.Stats
	fmt.Printf("%+v\n", stats)
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, wallets []wallet_mock.Wallet) {
	// TODO check all nodes are the same height ?
	// pick one node to draw height
	enclaveNode := l2Network.nodes[0].Enclave
	l1Height := enclaveNode.L1Height()
	l1HeightHash := enclaveNode.L1HeightHash()
	l2Height := enclaveNode.L2Height()
	l2HeightHash := enclaveNode.L2HeightHash()

	//stats := l1Network.Stats
	//l1BlockState, found := enclaveNode.TestDB().FetchState(l1Height)
	//if !found {
	//	fmt.Println("derp")
	//}
	validateL1(t, l1Network.Stats, l1Height, l1HeightHash, enclaveNode.TestDB())
	totalWithdrawn := validateL2(t, l1Network.Stats, enclaveNode, l2Height, l2HeightHash)
	validateL2State(t, l2Network, l1Network.Stats, totalWithdrawn, wallets)
}

// For this simulation, this represents an acceptable "dead blocks" percentage.
// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
// We test the results against this threshold to catch eventual protocol errors.
const L1EfficiencyThreashold = 0.2
const L2EfficiencyThreashold = 0.3

// Sanity check
func validateL1(t *testing.T, stats *Stats, l1Height int, l1HeightHash common.L1RootHash, db enclave.DB) {
	deposits := make([]uuid.UUID, 0)
	rollups := make([]common.L2RootHash, 0)
	stats.l1Height = l1Height
	totalDeposited := uint64(0)

	l1State, found := db.FetchState(l1HeightHash)
	if !found {
		t.Errorf("expected l1 height not found")
	}

	blockchain := ethereum_mock.BlocksBetween(&common.GenesisBlock, l1State.Block, db)
	headRollup := &enclave.GenesisRollup
	for _, block := range blockchain {
		for _, tx := range block.Transactions {
			currentRollups := make([]*enclave.Rollup, 0)
			switch tx.TxType {
			case common.DepositTx:
				deposits = append(deposits, tx.ID)
				totalDeposited += tx.Amount
			case common.RollupTx:
				r := obscuroCommon.DecodeRollup(tx.Rollup)
				rollups = append(rollups, r.Hash())
				if common.IsBlockAncestor(r.Header.L1Proof, block, db) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					currentRollups = append(currentRollups, enclave.DecryptRollup(r))
					stats.NewRollup(r)
				}
			default:
				panic("unknown transaction type")
			}
			r, _ := enclave.FindWinner(headRollup, currentRollups, db)
			if r != nil {
				headRollup = r
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

func validateL2(t *testing.T, stats *Stats, enclaveNode enclave.Enclave, l2Height int, l2HeightHash common.L2RootHash) uint64 {
	stats.l2Height = l2Height
	transfers := make([]uuid.UUID, 0)
	withdrawalTxs := make([]enclave.L2Tx, 0)
	withdrawalRequests := make([]obscuroCommon.Withdrawal, 0)

	// get transactions at current height (which we'll need to track somewhere else)
	transactions := enclaveNode.TransactionsAtHeight(l2HeightHash)
	withdrawals := enclaveNode.WithdrawlsAtHeight(l2HeightHash)

	for {
		if l2Height == common.L2GenesisHeight {
			break
		}
		for _, tx := range transactions {
			switch tx.TxType {
			case enclave.TransferTx:
				transfers = append(transfers, tx.ID)
			case enclave.WithdrawalTx:
				withdrawalTxs = append(withdrawalTxs, tx)
			default:
				panic("Invalid tx type")
			}
		}
		withdrawalRequests = append(withdrawalRequests, withdrawals...)
		l2Height--
		l2HeightHash = enclaveNode.ParentHash(l2HeightHash)
		transactions = enclaveNode.TransactionsAtHeight(l2HeightHash)
		withdrawals = enclaveNode.WithdrawlsAtHeight(l2HeightHash)
	}
	// todo - check that proofs are on the canonical chain

	if len(common.FindDups(transfers)) > 0 {
		dups := common.FindDups(transfers)
		t.Errorf("Found L2 txs duplicates: %v", dups)
	}
	if len(transfers) != stats.nrTransferTransactions {
		t.Errorf("Nr of transfers don't match. Found %d , expected %d", len(transfers), stats.nrTransferTransactions)
	}
	if sumWithdrawalTxs(withdrawalTxs) != stats.totalWithdrawnAmount {
		t.Errorf("Withdrawal tx amounts don't match. Found %d , expected %d", sumWithdrawalTxs(withdrawalTxs), stats.totalWithdrawnAmount)
	}
	if sumWithdrawals(withdrawalRequests) > stats.totalWithdrawnAmount {
		t.Errorf("The amount withdrawn %d exceeds the actual amount requested %d", sumWithdrawals(withdrawalRequests), stats.totalWithdrawnAmount)
	}
	efficiency := float64(stats.totalL2Blocks-stats.l2Height) / float64(stats.totalL2Blocks)
	if efficiency > L2EfficiencyThreashold {
		t.Errorf("Efficiency in L2 is %f. Expected:%f", efficiency, L2EfficiencyThreashold)
	}

	return sumWithdrawals(withdrawalRequests)
}

func sumWithdrawals(w []obscuroCommon.Withdrawal) uint64 {
	sum := uint64(0)
	for _, r := range w {
		sum += r.Amount
	}
	return sum
}

func sumWithdrawalTxs(t []enclave.L2Tx) uint64 {
	sum := uint64(0)
	for _, r := range t {
		sum += r.Amount
	}

	return sum
}

func validateL2State(t *testing.T, l2Network L2NetworkCfg, s *Stats, totalWithdrawn uint64, wallets []wallet_mock.Wallet) {
	finalAmount := s.totalDepositedAmount - totalWithdrawn

	// Parallelize this check
	var nGroup errgroup.Group

	// Check that the state on all nodes is valid
	for _, node := range l2Network.nodes {
		nGroup.Go(func() error {
			// add up all balances
			total := uint64(0)
			for _, wallet := range wallets {
				total += node.Enclave.Balance(wallet.Address)
			}
			if total != finalAmount {
				return fmt.Errorf("the amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", node.ID, total, finalAmount)
			}
			return nil
		})
	}
	err := nGroup.Wait()
	if err != nil {
		t.Error(err)
	}
	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// walk the blocks in reverse direction, execute deposits and transactions and compare to the state in the rollup
}

func totalBalance(s enclave.BlockState) uint64 {
	tot := uint64(0)
	for _, bal := range s.State {
		tot += bal
	}
	return tot
}
