package db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// this will eventually be the geth vm.StateDb
type obscuroStateDb struct {
	db         *inMemoryDB
	parentRoot obscurocommon.L2RootHash // root of the parent
	// currentRoot obscurocommon.L2RootHash
	state *State
}

func NewStateDb(db *inMemoryDB, root obscurocommon.L2RootHash, parentState *State) StateDb {
	return &obscuroStateDb{
		parentRoot: root,
		db:         db,
		state:      parentState,
	}
}

func (s *obscuroStateDb) GetBalance(address common.Address) uint64 {
	return s.state.Balances[address]
}

func (s *obscuroStateDb) SetBalance(address common.Address, balance uint64) {
	s.state.Balances[address] = balance
}

func (s *obscuroStateDb) AddWithdrawal(txHash obscurocommon.TxHash) {
	s.state.Withdrawals = append(s.state.Withdrawals, txHash)
}

func (s *obscuroStateDb) Commit(currentRoot obscurocommon.L2RootHash) {
	s.db.SetRollupState(currentRoot, s.state)
}

func (s *obscuroStateDb) Copy() StateDb {
	return NewStateDb(s.db, s.parentRoot, CopyState(s.state))
}

func (s *obscuroStateDb) StateRoot() common.Hash {
	return Serialize(s.state)
}

func (s *obscuroStateDb) Withdrawals() []obscurocommon.TxHash {
	return s.state.Withdrawals
}
