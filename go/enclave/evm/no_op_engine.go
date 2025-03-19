package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// NoOpConsensusEngine - implements the geth consensus.Engine, but doesn't do anything
// This is needed for running evm transactions
type NoOpConsensusEngine struct {
	logger gethlog.Logger
}

// Author is used to determine where to send the gas collected from the fees.
func (e *NoOpConsensusEngine) Author(_ *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

func (e *NoOpConsensusEngine) VerifyHeader(_ consensus.ChainHeaderReader, _ *types.Header) error {
	e.logger.Crit("noop")
	return nil
}

func (e *NoOpConsensusEngine) VerifyHeaders(_ consensus.ChainHeaderReader, _ []*types.Header) (chan<- struct{}, <-chan error) {
	e.logger.Crit("noop")
	return nil, nil
}

func (e *NoOpConsensusEngine) VerifyUncles(_ consensus.ChainReader, _ *types.Block) error {
	e.logger.Crit("noop")
	return nil
}

func (e *NoOpConsensusEngine) Prepare(_ consensus.ChainHeaderReader, _ *types.Header) error {
	e.logger.Crit("noop")
	return nil
}

func (e *NoOpConsensusEngine) Finalize(_ consensus.ChainHeaderReader, _ *types.Header, _ vm.StateDB, _ *types.Body) {
	panic("implement me")
}

func (e *NoOpConsensusEngine) FinalizeAndAssemble(_ consensus.ChainHeaderReader, _ *types.Header, _ *state.StateDB, _ *types.Body, _ []*types.Receipt) (*types.Block, error) {
	e.logger.Crit("noop")
	return nil, nil //nolint:nilnil
}

func (e *NoOpConsensusEngine) Seal(_ consensus.ChainHeaderReader, _ *types.Block, _ chan<- *types.Block, _ <-chan struct{}) error {
	e.logger.Crit("noop")
	return nil
}

func (e *NoOpConsensusEngine) SealHash(_ *types.Header) common.Hash {
	e.logger.Crit("noop")
	return common.Hash{}
}

func (e *NoOpConsensusEngine) CalcDifficulty(_ consensus.ChainHeaderReader, _ uint64, _ *types.Header) *big.Int {
	e.logger.Crit("noop")
	return nil
}

func (e *NoOpConsensusEngine) APIs(_ consensus.ChainHeaderReader) []rpc.API {
	e.logger.Crit("noop")
	return nil
}

func (e *NoOpConsensusEngine) Close() error {
	e.logger.Crit("noop")
	return nil
}
