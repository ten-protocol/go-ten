package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// NoOpEngine - implements the geth consensus.Engine, but doesn't do anything
// Todo - does it have to do anything?
type NoOpEngine struct{}

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (e *NoOpEngine) Author(header *types.Header) (common.Address, error) {
	h := convertFromEthHeader(header)
	return h.Agg, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (e *NoOpEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	panic("noop")
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (e *NoOpEngine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	panic("noop")
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (e *NoOpEngine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	panic("noop")
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (e *NoOpEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	panic("noop")
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// but does not assemble the block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (e *NoOpEngine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header,
) {
	panic("noop")
}

// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
// rewards) and assembles the final block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (e *NoOpEngine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt,
) (*types.Block, error) {
	panic("noop")
}

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (e *NoOpEngine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	panic("noop")
}

// SealHash returns the hash of a block prior to it being sealed.
func (e *NoOpEngine) SealHash(header *types.Header) common.Hash {
	panic("noop")
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have.
func (e *NoOpEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	panic("noop")
}

// APIs returns the RPC APIs this consensus engine provides.
func (e *NoOpEngine) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	panic("noop")
}

// Close terminates any background threads maintained by the consensus engine.
func (e *NoOpEngine) Close() error {
	panic("noop")
}
