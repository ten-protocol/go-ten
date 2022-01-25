package simulation

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"simulation/common"
	"testing"
	"time"
)

func TestSimulation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	// create a folder specific for the test
	d, err := os.MkdirTemp("..", "simulation_result")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(d + "/" + "simulation_result.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	common.SetLog(f)

	blockDuration := 15_000
	l1netw, l2netw := RunSimulation(5, 10, 20, blockDuration, blockDuration/20, blockDuration/3)
	checkBlockchainValidity(t, l1netw, l2netw)
}

// Checks that there are no duplicated transactions in the L1 or L2
//TODO check that all injected transactions were included
//TODO Check that the total amount of money in user accounts matches the amount injected as deposits
//TODO Check that processing transactions in the order specified in the list results in the same balances
func checkBlockchainValidity(t *testing.T, l1Network L1NetworkCfg, l2Network L2NetworkCfg) {
	r := l1Network.Stats.l2Head

	// check that there are no duplicate transactions on the L1
	deposits := make([]uuid.UUID, 0)
	rollups := make([]uuid.UUID, 0)
	b := r.L1Proof
	totalTx := 0

	for {
		if b.Height() == -1 {
			break
		}
		for _, tx := range b.L1Txs() {
			if tx.TxType == common.DepositTx {
				deposits = append(deposits, tx.Id)
				totalTx += tx.Amount
			} else {
				rollups = append(rollups, tx.Rollup.RootHash())
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

	transfers := make([]uuid.UUID, 0)
	for {
		if r.Height() == -1 {
			break
		}
		for _, tx := range r.L2Txs() {
			if tx.TxType == common.TransferTx {
				transfers = append(transfers, tx.Id)
				//totalTx += tx.Amount
			}
		}
		r = r.ParentRollup()
	}

	if len(findDups(transfers)) > 0 {
		dups := findDups(transfers)
		t.Errorf("Found L2 txs duplicates: %v", dups)
	}

	//fmt.Printf("Deposits: total_in=%d; total_txs=%d\n", total, totalTx)

	bl := l1Network.Stats.l2Head.L1Proof

	nrDeposits := 0
	totalDeposits := 0

	// walk the L1 blocks and
	for {

		if bl.Height() == -1 {
			break
		}

		s, _ := l2Network.nodes[0].Db.Fetch(bl.RootHash())
		tot := 0
		for _, bal := range s.State {
			tot += bal
		}
		nrDeposits = 0
		for _, tx := range bl.L1Txs() {
			if tx.TxType == common.DepositTx {
				nrDeposits++
			}
		}
		totalDeposits += nrDeposits

		fmt.Printf("%d=%d (%d of %d)\n", bl.RootHash().ID(), tot, nrDeposits, totalDeposits)
		//lastTotal = t
		bl = bl.ParentBlock()
	}

}
