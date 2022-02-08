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

	blockDuration := uint64(30_000)
	l1netw, l2netw := RunSimulation(5, 30, 30, blockDuration, blockDuration/15, blockDuration/4)
	stats := l1netw.Stats
	fmt.Printf("%+v\n", stats)
	firstNode := l2netw.nodes[0]
	checkBlockchainValidity(t, l1netw, l2netw, firstNode.Enclave.Db(), firstNode.Enclave.PeekHead().Head)
	//pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, db obscuro.Db, r obscuro.Rollup) {
	stats := l1Network.Stats
	p := r.Proof(db)
	validateL1(t, p, stats, db)
	totalWithdrawn := validateL2(t, r, stats, db)
	validateL2State(t, l1Network, l2Network, stats, totalWithdrawn)
}

// For this simulation, this represents an acceptable "dead blocks" percentage.
// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
// We test the results against this threshold to catch eventual protocol errors.
const L1EfficiencyThreashold = 0.2
const L2EfficiencyThreashold = 0.3

// Sanity check
func validateL1(t *testing.T, b common.Block, s *Stats, db obscuro.Db) {
	deposits := make([]uuid.UUID, 0)
	rollups := make([]uuid.UUID, 0)

	totalDeposited := uint64(0)

	for {
		if b.Height == common.GenesisHeight {
			break
		}
		for _, tx := range b.Transactions {
			switch tx.TxType {
			case common.DepositTx:
				deposits = append(deposits, tx.Id)
				totalDeposited += tx.Amount
			case common.RollupTx:
				rollups = append(rollups, obscuro.DecodeRollup(tx.Rollup).RootHash)
			default:
				panic("unknown transaction type")
			}
		}
		p, f := b.Parent(db)
		if !f {
			panic("wtf")
		}

		b = p
	}

	if len(common.FindDups(deposits)) > 0 {
		dups := common.FindDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(common.FindDups(rollups)) > 0 {
		dups := common.FindDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != s.totalDepositedAmount {
		t.Errorf("Deposit amounts don't match. Found %d , expected %d", totalDeposited, s.totalDepositedAmount)
	}

	efficiency := float64(uint32(s.totalL1Blocks)-s.l1Height) / float64(s.totalL1Blocks)
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

func validateL2(t *testing.T, r obscuro.Rollup, s *Stats, db obscuro.Db) uint64 {
	transfers := make([]uuid.UUID, 0)
	withdrawalTxs := make([]obscuro.L2Tx, 0)
	withdrawalRequests := make([]obscuro.Withdrawal, 0)
	for {
		if r.Height == common.GenesisHeight {
			break
		}
		for _, tx := range r.Transactions {
			switch tx.TxType {
			case obscuro.TransferTx:
				transfers = append(transfers, tx.Id)
			case obscuro.WithdrawalTx:
				withdrawalTxs = append(withdrawalTxs, tx)
			default:
				panic("Invalid tx type")
			}
		}
		withdrawalRequests = append(withdrawalRequests, r.Withdrawals...)
		r = r.Parent(db)
	}
	//todo - check that proofs are on the canonical chain

	if len(common.FindDups(transfers)) > 0 {
		dups := common.FindDups(transfers)
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
	efficiency := float64(uint32(s.totalL2Blocks)-s.l2Height) / float64(s.totalL2Blocks)
	if efficiency > L2EfficiencyThreashold {
		t.Errorf("Efficiency in L2 is %f. Expected:%f", efficiency, L1EfficiencyThreashold)
	}

	return sumWithdrawals(withdrawalRequests)
}

func sumWithdrawals(w []obscuro.Withdrawal) uint64 {
	sum := uint64(0)
	for _, r := range w {
		sum += r.Amount
	}
	return sum
}

func sumWithdrawalTxs(t []obscuro.L2Tx) uint64 {
	sum := uint64(0)
	for _, r := range t {
		sum += r.Amount
	}
	return sum
}

func validateL2State(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, s *Stats, totalWithdrawn uint64) {

	finalAmount := s.totalDepositedAmount - totalWithdrawn
	// Check that the state on all nodes is valid
	for _, observer := range l2Network.nodes {
		// read the last state
		lastState := observer.Enclave.PeekHead()
		total := totalBalance(lastState)
		if total != finalAmount {
			t.Errorf("The amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", observer.Id, total, finalAmount)
		}
	}

	//TODO Check that processing transactions in the order specified in the list results in the same balances
	// walk the blocks in reverse direction, execute deposits and transactions and compare to the state in the rollup
}

func totalBalance(s obscuro.BlockState) uint64 {
	tot := uint64(0)
	for _, bal := range s.State {
		tot += bal
	}
	return tot
}
