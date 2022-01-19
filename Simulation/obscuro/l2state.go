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

type Db interface {
	fetch(hash L1RootHash) (BlockState, bool)
	set(hash L1RootHash, state BlockState)
}

type InMemoryDb struct {
	// the state is dependent on the L1 block alone
	cache map[L1RootHash]BlockState
	mutex sync.RWMutex
}

func NewInMemoryDb() *InMemoryDb {
	return &InMemoryDb{
		cache: make(map[L1RootHash]BlockState),
		mutex: sync.RWMutex{},
	}
}

func (db *InMemoryDb) fetch(hash L1RootHash) (BlockState, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	val, found := db.cache[hash]
	return val, found
}

func (db *InMemoryDb) set(hash L1RootHash, state BlockState) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.cache[hash] = state
}

func copyState(state State) State {
	s := make(State)
	for address, balance := range state {
		s[address] = balance
	}
	return s
}

func serialize(state State) string {
	s := make([]string, 0)
	for add, bal := range state {
		s = append(s, fmt.Sprintf("%d=%d", add.ID(), bal))
	}
	return fmt.Sprintf("%v", s)
}

// returns a modified copy of the state
func calculateState(txs []*L2Tx, state State) State {
	s := copyState(state)
	for _, tx := range txs {
		executeTx(s, tx)
	}
	return s
}

// mutates the state
func executeTx(s State, tx *L2Tx) {
	bal, _ := s[tx.from]
	if bal >= tx.amount {
		s[tx.from] -= tx.amount
		s[tx.dest] += tx.amount
	}
}

type RollupStorage struct {
}

func emptyState() State {
	return make(State)
}
