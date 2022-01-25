package obscuro

import (
	"fmt"
	"simulation/common"
	"simulation/wallet-mock"
	"sync"
)

type State = map[wallet_mock.Address]int

type BlockState struct {
	Head  *common.Rollup
	State State
}

type Db interface {
	Fetch(hash common.RootHash) (BlockState, bool)
	Set(hash common.RootHash, state BlockState)
}

type InMemoryDb struct {
	// the State is dependent on the L1 block alone
	cache map[common.RootHash]BlockState
	mutex sync.RWMutex
}

func NewInMemoryDb() *InMemoryDb {
	return &InMemoryDb{
		cache: make(map[common.RootHash]BlockState),
		mutex: sync.RWMutex{},
	}
}

func (db *InMemoryDb) Fetch(hash common.RootHash) (BlockState, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	val, found := db.cache[hash]
	return val, found
}

func (db *InMemoryDb) Set(hash common.RootHash, state BlockState) {
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
	return fmt.Sprintf("%v", state)
}

// returns a modified copy of the State
func calculateState(txs []common.L2Tx, state State) State {
	s := copyState(state)
	for _, tx := range txs {
		executeTx(s, tx)
	}
	return s
}

// mutates the State
func executeTx(s State, tx common.L2Tx) {
	bal, _ := s[tx.From]
	if bal >= tx.Amount {
		s[tx.From] -= tx.Amount
		s[tx.Dest] += tx.Amount
		//} else {
		//fmt.Printf("--%d\n", tx.Id.Id())
	}
}

type RollupStorage struct {
}

func emptyState() State {
	return make(State)
}
