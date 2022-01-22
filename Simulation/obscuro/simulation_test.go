package obscuro

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestSimulation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	// todo - create a folder specific for the test
	d, err := os.MkdirTemp("..", "simulation_result")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(d + "/" + "simulation_result.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	SetLog(f)

	blockDuration := 15_000
	netw := RunSimulation(5, 10, 20, blockDuration, blockDuration/20, blockDuration/3)
	checkBlockchainValidity(netw, t)
}

// Checks that there are no duplicated transactions in the L1 or L2
//TODO check that all injected transactions were included
//TODO Check that the total amount of money in user accounts matches the amount injected as deposits
//TODO Check that processing transactions in the order specified in the list results in the same balances
func checkBlockchainValidity(network NetworkCfg, t *testing.T) {
	r := network.Stats.l2Head

	// check that there are no duplicate transactions on the L1
	deposits := make([]uuid.UUID, 0)
	rollups := make([]uuid.UUID, 0)
	b := r.l1Proof
	totalTx := 0

	for {
		if b.rootHash == GenesisBlock.rootHash {
			break
		}
		for _, tx := range b.txs {
			if tx.txType == DepositTx {
				deposits = append(deposits, tx.id)
				totalTx += tx.amount
			} else {
				rollups = append(rollups, tx.rollup.rootHash)
			}
		}
		b = b.parent
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
		if r.rootHash == GenesisRollup.rootHash {
			break
		}
		for _, tx := range r.txs {
			if tx.txType == TransferTx {
				transfers = append(transfers, tx.id)
				//totalTx += tx.amount
			}
		}
		r = r.parent
	}

	if len(findDups(transfers)) > 0 {
		dups := findDups(transfers)
		t.Errorf("Found L2 txs duplicates: %v", dups)
	}

	//fmt.Printf("Deposits: total_in=%d; total_txs=%d\n", total, totalTx)

	bl := network.Stats.l2Head.l1Proof

	nrDeposits := 0
	totalDeposits := 0

	// walk the L1 blocks and
	for {

		if bl.rootHash == GenesisBlock.rootHash {
			break
		}

		s, _ := network.allAgg[0].db.fetch(bl.rootHash)
		tot := 0
		for _, bal := range s.state {
			tot += bal
		}
		nrDeposits = 0
		for _, tx := range bl.txs {
			if tx.txType == DepositTx {
				nrDeposits++
			}
		}
		totalDeposits += nrDeposits

		fmt.Printf("%d=%d (%d of %d)\n", bl.rootHash.ID(), tot, nrDeposits, totalDeposits)
		//lastTotal = t
		bl = bl.parent
	}

}
