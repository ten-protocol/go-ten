package genesis

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/holiman/uint256"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ten-protocol/go-ten/go/common"
)

// Account specifies the address that's prefunded and the amount it's funded with
type Account struct {
	Address gethcommon.Address `json:"address"`
	Amount  *big.Int           `json:"amount"`
}

// Contract specifies an address and its bytecode to be set at genesis
type Contract struct {
	Address  gethcommon.Address `json:"address"`
	Bytecode string             `json:"bytecode"` // hex string with 0x prefix
}

// Genesis holds a range of prefunded accounts
type Genesis struct {
	Accounts  []Account  `json:"accounts"`
	Contracts []Contract `json:"contracts"`
}

// New creates a new Genesis given a json string
// if the string is empty it defaults to the testnet genesis
func New(genesisJSON string) (*Genesis, error) {
	genesis := &Genesis{}
	if len(genesisJSON) == 0 {
		return genesis, nil
	}
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
	_, err = stateDB.Commit(0, false, true)
	if err != nil {
		return err
	}

	// todo - VERKLE
	//if root != (gethcommon.Hash{}) {
	//	if err := storage.TrieDB().Commit(root, true); err != nil {
	//		return err
	//	}
	//}

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

	// set predeployed contract code
	for _, c := range g.Contracts {
		if c.Bytecode == "" {
			continue
		}
		code, err := hexutil.Decode(c.Bytecode)
		if err != nil {
			return nil, fmt.Errorf("invalid contract bytecode for %s: %w", c.Address.Hex(), err)
		}
		s.SetCode(c.Address, code)
	}

	return s, nil
}
