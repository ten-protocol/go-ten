package obscuro

import (
	"fmt"
	"simulation/common"
	"simulation/wallet-mock"
	"sync"
)

type State = map[wallet_mock.Address]int

type Db interface {
	Fetch(hash common.RootHash) (BlockState, bool)
	Set(hash common.RootHash, state BlockState)
	Head() BlockState
	Balance(address wallet_mock.Address) int
}

type InMemoryDb struct {
	// the State is dependent on the L1 block alone
	cache map[common.RootHash]BlockState
	head  common.RootHash
	mutex sync.RWMutex
}

// BlockState - Represents the state after an L1 block was processed.
type BlockState struct {
	Block          common.Block
	Head           common.Rollup
	State          State
	foundNewRollup bool
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

	// todo - is there any other logic here?
	db.head = hash
}

func (db *InMemoryDb) Head() BlockState {
	val, _ := db.Fetch(db.head)
	return val
}

func (db *InMemoryDb) Balance(address wallet_mock.Address) int {
	return db.Head().State[address]
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
func executeTransactions(txs []common.L2Tx, state State) State {
	s := copyState(state)
	for _, tx := range txs {
		executeTx(s, tx)
	}
	return s
}

// mutates the State
func executeTx(s State, tx common.L2Tx) {
	bal := s[tx.From]
	if bal >= tx.Amount {
		s[tx.From] -= tx.Amount
		s[tx.Dest] += tx.Amount
	}
}

func emptyState() State {
	return make(State)
}
