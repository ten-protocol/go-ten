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
// TODO - this has to be reworked when the gas/fee work starts
var PoolAddress = common.HexToAddress("0x0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A0A")

// ObscuroNoOpConsensusEngine - implements the geth consensus.Engine, but doesn't do anything
// This is needed for running evm transactions
type ObscuroNoOpConsensusEngine struct {
	logger gethlog.Logger
}

// Author is used to determine where to send the gas collected from the fees.
func (e *ObscuroNoOpConsensusEngine) Author(header *types.Header) (common.Address, error) {
	return PoolAddress, nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	e.logger.Crit("noop")
	return nil, nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	e.logger.Crit("noop")
}

func (e *ObscuroNoOpConsensusEngine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt,
) (*types.Block, error) {
	e.logger.Crit("noop")
	return nil, nil // nolint:nilnil
}

func (e *ObscuroNoOpConsensusEngine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) SealHash(header *types.Header) common.Hash {
	e.logger.Crit("noop")
	return common.Hash{}
}

func (e *ObscuroNoOpConsensusEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	e.logger.Crit("noop")
	return nil
}

func (e *ObscuroNoOpConsensusEngine) Close() error {
	e.logger.Crit("noop")
	return nil
}
