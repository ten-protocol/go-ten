package obscuro

import (
	"fmt"
	"sync"
)

type State = map[Address]int

type BlockState struct {
	head  *Rollup
	state State
}

// the state is dependent on the L1 block alone
var globalDb = make(map[L1RootHash]BlockState)
var dbMutex = &sync.RWMutex{}

func copyState(state State) State {
	s := make(State)
	for address, balance := range state {
		s[address] = balance
	}
	return s
}

func serialize(state State) StateRoot {
	s := make([]string, 0)
	for add, bal := range state {
		s = append(s, fmt.Sprintf("%d=%d", add.ID(), bal))
	}
	return fmt.Sprintf("%v", s)
}

// returns a modified copy of the state
func (a L2Agg) calculateState(txs []*L2Tx, state State) State {
	s := copyState(state)
	for _, tx := range txs {
		executeTx(s, tx)
	}
	return s
}

func executeTx(s State, tx *L2Tx) {
	bal, _ := s[tx.from]
	if bal >= tx.amount {
		s[tx.from] -= tx.amount
		s[tx.dest] += tx.amount
	}
}

type RollupStorage struct {
}
