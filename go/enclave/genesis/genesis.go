package genesis

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/holiman/uint256"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ten-protocol/go-ten/go/common"
)

// Account specifies the address that's prefunded and the amount it's funded with
type Account struct {
	Address gethcommon.Address
	Amount  *big.Int
}

// Genesis holds a range of prefunded accounts
type Genesis struct {
	Accounts []Account
}

// New creates a new Genesis given a json string
// if the string is empty it defaults to the testnet genesis
func New(genesisJSON string) (*Genesis, error) {
	genesis := &Genesis{}
	err := json.Unmarshal([]byte(genesisJSON), genesis)
	if err != nil {
		return nil, err
	}
	return genesis, nil
}

func (g Genesis) CommitGenesisState(storage storage.Storage) error {
	stateDB, err := g.applyAllocations(storage)
	if err != nil {
		return err
	}
	_, err = stateDB.Commit(0, false)
	if err != nil {
		return err
	}
	return nil
}

func (g Genesis) GetGenesisRoot(storage storage.Storage) (*common.StateRoot, error) {
	stateDB, err := g.applyAllocations(storage)
	if err != nil {
		return nil, err
	}
	stateHash := stateDB.IntermediateRoot(true)
	return &stateHash, nil
}

// Applies the faucet preallocation on top of an empty state DB.
func (g Genesis) applyAllocations(storage storage.Storage) (*state.StateDB, error) {
	s, err := storage.EmptyStateDB()
	if err != nil {
		return nil, fmt.Errorf("could not initialise empty state DB. Cause: %w", err)
	}

	// set the accounts funds
	for _, acc := range g.Accounts {
		s.SetBalance(acc.Address, uint256.MustFromBig(acc.Amount), tracing.BalanceIncreaseGenesisBalance)
	}

	return s, nil
}
