package simulation

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	common2 "github.com/ethereum/go-ethereum/common"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/common"
	obscuroCommon "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
	enclave2 "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func TestSimulation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	// create a folder specific for the test
	err := os.MkdirAll("../.build/simulations/", 0o700)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("simulation-result-%d-*.txt", time.Now().Unix())
	f, err := os.CreateTemp("../.build/simulations", fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	common.SetLog(f)

	blockDuration := uint64(20_000)
	l1netw, l2netw := RunSimulation(5, 10, 15, blockDuration, blockDuration/15, blockDuration/3)
	firstNode := l2netw.nodes[0]
	checkBlockchainValidity(t, l1netw, l2netw, firstNode.Enclave.TestDB(), firstNode.Enclave.TestPeekHead().Head)
	stats := l1netw.Stats
	fmt.Printf("%+v\n", stats)
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg, db enclave2.DB, r *enclave2.Rollup) {
	stats := l1Network.Stats
	p := r.Proof(db)
	validateL1(t, p, stats, db)
	totalWithdrawn := validateL2(t, r, stats, db)
	validateL2State(t, l2Network, stats, totalWithdrawn)
}

// For this simulation, this represents an acceptable "dead blocks" percentage.
// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
// We test the results against this threshold to catch eventual protocol errors.
const L1EfficiencyThreashold = 0.2
const L2EfficiencyThreashold = 0.3

// Sanity check
func validateL1(t *testing.T, b *common.Block, s *Stats, db enclave2.DB) {
	deposits := make([]uuid.UUID, 0)
	rollups := make([]common.L2RootHash, 0)
	s.l1Height = b.Height(db)
	totalDeposited := uint64(0)

	blockchain := ethereum_mock.BlocksBetween(&common.GenesisBlock, b, db)
	headRollup := &enclave2.GenesisRollup
	for _, block := range blockchain {
		for _, tx := range block.Transactions {
			currentRollups := make([]*enclave2.Rollup, 0)
			switch tx.TxType {
			case common.DepositTx:
				deposits = append(deposits, tx.ID)
				totalDeposited += tx.Amount
			case common.RollupTx:
				r := obscuroCommon.DecodeRollup(tx.Rollup)
				rollups = append(rollups, r.Hash())
				if common.IsBlockAncestor(r.Header.L1Proof, b, db) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					currentRollups = append(currentRollups, enclave2.DecryptRollup(r))
					s.NewRollup(r)
				}
			case common.RequestSecretTx:
			case common.StoreSecretTx:
			}
			r, _ := enclave2.FindWinner(headRollup, currentRollups, db)
			if r != nil {
				headRollup = r
			}
		}
	}

	if len(common.FindUUIDDups(deposits)) > 0 {
		dups := common.FindUUIDDups(deposits)
		t.Errorf("Found Deposit duplicates: %v", dups)
	}
	if len(common.FindRollupDups(rollups)) > 0 {
		dups := common.FindRollupDups(rollups)
		t.Errorf("Found Rollup duplicates: %v", dups)
	}
	if totalDeposited != s.totalDepositedAmount {
		t.Errorf("Deposit amounts don't match. Found %d , expected %d", totalDeposited, s.totalDepositedAmount)
	}

	efficiency := float64(s.totalL1Blocks-s.l1Height) / float64(s.totalL1Blocks)
	if efficiency > L1EfficiencyThreashold {
		t.Errorf("Efficiency in L1 is %f. Expected:%f", efficiency, L1EfficiencyThreashold)
	}

	// todo
	for nodeID, reorgs := range s.noL1Reorgs {
		eff := float64(reorgs) / float64(s.l1Height)
		if eff > L1EfficiencyThreashold {
			t.Errorf("Efficiency for node %d in L1 is %f. Expected:%f", nodeID, eff, L1EfficiencyThreashold)
		}
	}
}

func validateL2(t *testing.T, r *enclave2.Rollup, s *Stats, db enclave2.DB) uint64 {
	s.l2Height = db.Height(r)
	transfers := make([]common2.Hash, 0)
	withdrawalTxs := make([]enclave2.L2Tx, 0)
	withdrawalRequests := make([]obscuroCommon.Withdrawal, 0)
	for {
		if db.Height(r) == common.L2GenesisHeight {
			break
		}
		for i := range r.Transactions {
			tx := r.Transactions[i]
			txData := enclave2.TxData(&tx)
			switch txData.Type {
			case enclave2.TransferTx:
				transfers = append(transfers, tx.Hash())
			case enclave2.WithdrawalTx:
				withdrawalTxs = append(withdrawalTxs, tx)
			default:
				panic("Invalid tx type")
			}
		}
		withdrawalRequests = append(withdrawalRequests, r.Header.Withdrawals...)
		r = db.Parent(r)
	}
	// todo - check that proofs are on the canonical chain

	if len(common.FindHashDups(transfers)) > 0 {
		dups := common.FindHashDups(transfers)
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
	efficiency := float64(s.totalL2Blocks-s.l2Height) / float64(s.totalL2Blocks)
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

func sumWithdrawalTxs(t []enclave2.L2Tx) uint64 {
	sum := uint64(0)
	for i := range t {
		txData := enclave2.TxData(&t[i])
		sum += txData.Amount
	}

	return sum
}

func validateL2State(t *testing.T, l2Network L2NetworkCfg, s *Stats, totalWithdrawn uint64) {
	finalAmount := s.totalDepositedAmount - totalWithdrawn
	// Check that the state on all nodes is valid
	for _, observer := range l2Network.nodes {
		// read the last state
		lastState := observer.Enclave.TestPeekHead()
		total := totalBalance(lastState)
		if total != finalAmount {
			t.Errorf("The amount of money in accounts on node %d does not match the amount deposited. Found %d , expected %d", observer.ID, total, finalAmount)
		}
	}

	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// walk the blocks in reverse direction, execute deposits and transactions and compare to the state in the rollup
}

func totalBalance(s enclave2.BlockState) uint64 {
	tot := uint64(0)
	for _, bal := range s.State {
		tot += bal
	}
	return tot
}
