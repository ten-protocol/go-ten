package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ObscuroNoOpConsensusEngine - implements the geth consensus.Engine, but doesn't do anything
// This is needed for running evm transactions
type ObscuroNoOpConsensusEngine struct{}

func (e *ObscuroNoOpConsensusEngine) Author(header *types.Header) (common.Address, error) {
	h := convertFromEthHeader(header)
	return h.Agg, nil
}

func (e *ObscuroNoOpConsensusEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header,
) {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt,
) (*types.Block, error) {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) SealHash(header *types.Header) common.Hash {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	panic("noop")
}

func (e *ObscuroNoOpConsensusEngine) Close() error {
	panic("noop")
}
