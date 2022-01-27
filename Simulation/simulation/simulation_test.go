package simulation

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"simulation/common"
	"simulation/obscuro"
	"testing"
	"time"
)

func TestSimulation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	// create a folder specific for the test
	err := os.MkdirAll("../.build/simulations/", 0700)
	if err != nil {
		panic(err)
	}
	f, err := os.CreateTemp("../.build/simulations", "simulation-result-*.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	common.SetLog(f)

	blockDuration := 20_000
	l1netw, l2netw := RunSimulation(5, 10, 30, blockDuration, blockDuration/20, blockDuration/4)
	checkBlockchainValidity(t, l1netw, l2netw)
}

func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg) {
	r := l1Network.Stats.l2Head
	stats := l1Network.Stats
	fmt.Printf("%#v\n", stats)

	validateL1(t, r.L1Proof, stats)
	totalWithdrawn := validateL2(t, r, stats)
	validateL2State(t, l1Network, l2Network, stats, totalWithdrawn)
}

// For this simulation, this represents an acceptable "dead blocks" percentage.
// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
// We test the results against this threshold to catch eventual protocol errors.
const L1EfficiencyThreashold = 0.2

// Sanity check
func validateL1(t *testing.T, b *common.Block, s *Stats) {
	deposits := make([]uuid.UUID, 0)
	rollups := make([]uuid.UUID, 0)

	totalDeposited := 0

	for {
		if b.Height() == common.GenesisHeight {
			break
		}
		for _, tx := range b.L1Txs() {
			switch tx.TxType {
			case common.DepositTx:
				deposits = append(deposits, tx.Id)
				totalDeposited += tx.Amount
			case common.RollupTx:
				rollups = append(rollups, tx.Rollup.RootHash())
			default:
				panic("unknown transaction type")
			}
		}
		b = b.ParentBlock()
	}

	if len(findDups(deposits)) > 0 {
		dups := findDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(findDups(rollups)) > 0 {
		dups := findDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != s.totalDepositedAmount {
		t.Errorf("Deposit amounts don't match. Found %d , expected %d", totalDeposited, s.totalDepositedAmount)
	}

	efficiency := float64(s.totalL1Blocks-s.l1Height) / float64(s.totalL1Blocks)
	if efficiency > L1EfficiencyThreashold {
		t.Errorf("Efficiency in L1 is %f. Expected:%f", efficiency, L1EfficiencyThreashold)
	}

	//todo
	for nodeId, reorgs := range s.noL1Reorgs {
		eff := float64(reorgs) / float64(s.l1Height)
		if eff > L1EfficiencyThreashold {
			t.Errorf("Efficiency for node %d in L1 is %f. Expected:%f", nodeId, eff, L1EfficiencyThreashold)
		}
	}
}

func validateL2(t *testing.T, r *common.Rollup, s *Stats) int {
	transfers := make([]uuid.UUID, 0)
	withdrawalTxs := make([]common.L2Tx, 0)
	withdrawalRequests := make([]common.Withdrawal, 0)
	for {
		if r.Height() == common.GenesisHeight {
			break
		}
		for _, tx := range r.L2Txs() {
			switch tx.TxType {
			case common.TransferTx:
				transfers = append(transfers, tx.Id)
			case common.WithdrawalTx:
				withdrawalTxs = append(withdrawalTxs, tx)
			default:
				panic("Invalid tx type")
			}
		}
		withdrawalRequests = append(withdrawalRequests, r.Withdrawals...)
		r = r.ParentRollup()
	}
	//todo - check that proofs are on the canonical chain

	if len(findDups(transfers)) > 0 {
		dups := findDups(transfers)
		t.Errorf("Found L2 txs duplicates: %v", dups)
	}
	if len(transfers) != s.nrTransferTransactions {
		t.Errorf("Nr of transfers don't match. Found %d , expected %d", len(transfers), s.nrTransferTransactions)
	}
	if sumWithdrawalTxs(withdrawalTxs) != s.totalWithdrawnAmount {
		t.Errorf("Withdrawal tx amounts don't match. Found %d , expected %d", sumWithdrawalTxs(withdrawalTxs), s.totalWithdrawnAmount)
	}
	if sumWithdrawals(withdrawalRequests) > s.totalWithdrawnAmount {
		t.Errorf("The amount withdrawn %d exceeds the actual amount requested %d", sumWithdrawals(withdrawalRequests), s.totalWithdrawnAmount)
	}
	return sumWithdrawals(withdrawalRequests)
}

func sumWithdrawals(w []common.Withdrawal) int {
	sum := 0
	for _, r := range w {
		sum += r.Amount
	}
	return sum
}

func sumWithdrawalTxs(t []common.L2Tx) int {
	sum := 0
	for _, r := range t {
		sum += r.Amount
	}
	return sum
}

func validateL2State(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, s *Stats, totalWithdrawn int) {

	finalAmount := s.totalDepositedAmount - totalWithdrawn
	// Check that the state on all nodes is valid
	for _, observer := range l2Network.nodes {
		// read the last state
		lastState := observer.Db.Head()
		total := totalBalance(lastState)
		if total != finalAmount {
			t.Errorf("The amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", observer.Id, total, finalAmount)
		}
	}

	//TODO Check that processing transactions in the order specified in the list results in the same balances
	// walk the blocks in reverse direction, execute deposits and transactions and compare to the state in the rollup
}

func totalBalance(s obscuro.BlockState) int {
	tot := 0
	for _, bal := range s.State {
		tot += bal
	}
	return tot
}
