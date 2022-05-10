package db

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/hashing"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// State - this is a placeholder for the real Trie based state
//- people send transactions to an ObsERC20 that was a withdraw(amount, from, destination) method
//In the EVM, there will be a smart contract that does the following:
//- the tokens are deducted from the "from" address , and burned
//- add to the "Withdrawals" transactions - this info will be taken from the state
//Post processing, outside the evm:
//- generate withdrawal instructions (amount, destination), based on which withdrawal transaction were executed successfully
type State struct {
	Balances    map[common.Address]uint64
	Withdrawals []obscurocommon.TxHash
}

func CopyStateNoWithdrawals(state *State) *State {
	s := EmptyState()
	if state == nil {
		return s
	}
	for address, balance := range state.Balances {
		s.Balances[address] = balance
	}
	return s
}

func CopyState(state *State) *State {
	s := EmptyState()
	if state == nil {
		return s
	}
	for address, balance := range state.Balances {
		s.Balances[address] = balance
	}
	s.Withdrawals = append(s.Withdrawals, state.Withdrawals...)
	return s
}

func Serialize(state *State) nodecommon.StateRoot {
	hash, err := hashing.RLPHash(fmt.Sprintf("%v", state))
	if err != nil {
		panic(err)
	}
	return hash
}

func EmptyState() *State {
	return &State{
		Balances:    map[common.Address]uint64{},
		Withdrawals: make([]obscurocommon.TxHash, 0),
	}
}
