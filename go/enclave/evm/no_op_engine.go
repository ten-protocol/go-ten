package evm

import (
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// PoolAddress - address of the pool where the gas will go
// todo (#627) - this has to be reworked when the gas/fee work starts
var PoolAddress = common.HexToAddress("0x0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A")

// ObscuroNoOpConsensusEngine - implements the geth consensus.Engine, but doesn't do anything
// This is needed for running evm transactions
type ObscuroNoOpConsensusEngine struct {
	logger gethlog.Logger
}

// Author is used to determine where to send the gas collected from the fees.
func (e *ObscuroNoOpConsensusEngine) Author(_ *types.Header) (common.Address, error) {
	return PoolAddress, nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyHeader(_ consensus.ChainHeaderReader, _ *types.Header) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyHeaders(_ consensus.ChainHeaderReader, _ []*types.Header) (chan<- struct{}, <-chan error) {
	e.logger.Crit("noop")
	return nil, nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyUncles(_ consensus.ChainReader, _ *types.Block) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) Prepare(_ consensus.ChainHeaderReader, _ *types.Header) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) Finalize(_ consensus.ChainHeaderReader, _ *types.Header, _ *state.StateDB, _ *types.Body) {
	// TODO implement me
	panic("implement me")
}

func (e *ObscuroNoOpConsensusEngine) FinalizeAndAssemble(_ consensus.ChainHeaderReader, _ *types.Header, _ *state.StateDB, _ *types.Body, _ []*types.Receipt) (*types.Block, error) {
	e.logger.Crit("noop")
	return nil, nil //nolint:nilnil
}

func (e *ObscuroNoOpConsensusEngine) Seal(_ consensus.ChainHeaderReader, _ *types.Block, _ chan<- *types.Block, _ <-chan struct{}) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) SealHash(_ *types.Header) common.Hash {
	e.logger.Crit("noop")
	return common.Hash{}
}

func (e *ObscuroNoOpConsensusEngine) CalcDifficulty(_ consensus.ChainHeaderReader, _ uint64, _ *types.Header) *big.Int {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) APIs(_ consensus.ChainHeaderReader) []rpc.API {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) Close() error {
	e.logger.Crit("noop")
	return nil
}
