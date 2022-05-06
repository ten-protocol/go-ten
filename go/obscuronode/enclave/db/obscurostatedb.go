package db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// this will eventually be the geth vm.StateDB
type obscuroStateDB struct {
	db         *InMemoryDB
	parentRoot obscurocommon.L2RootHash // root of the parent
	// currentRoot obscurocommon.L2RootHash
	state *State
}

func NewStateDB(db *InMemoryDB, root obscurocommon.L2RootHash, parentState *State) StateDB {
	return &obscuroStateDB{
		parentRoot: root,
		db:         db,
		state:      parentState,
	}
}

func (s *obscuroStateDB) GetBalance(address common.Address) uint64 {
	return s.state.Balances[address]
}

func (s *obscuroStateDB) SetBalance(address common.Address, balance uint64) {
	s.state.Balances[address] = balance
}

func (s *obscuroStateDB) AddWithdrawal(txHash obscurocommon.TxHash) {
	s.state.Withdrawals = append(s.state.Withdrawals, txHash)
}

func (s *obscuroStateDB) Commit(currentRoot obscurocommon.L2RootHash) {
	s.db.SetRollupState(currentRoot, s.state)
}

func (s *obscuroStateDB) Copy() StateDB {
	return NewStateDB(s.db, s.parentRoot, CopyState(s.state))
}

func (s *obscuroStateDB) StateRoot() common.Hash {
	return Serialize(s.state)
}

func (s *obscuroStateDB) Withdrawals() []obscurocommon.TxHash {
	return s.state.Withdrawals
}
