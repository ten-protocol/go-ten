package common

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var GenesisHash = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")

// GenesisBlock - this is a hack that has to be removed ASAP.
var GenesisBlock = NewBlock(nil, common.HexToAddress("0x0"), []*types.Transaction{})

// EncodedBlock the encoded version of an ExtBlock
type EncodedBlock []byte

func EncodeBlock(b *types.Block) (EncodedBlock, error) {
	encoded, err := rlp.EncodeToBytes(b)
	if err != nil {
		return nil, fmt.Errorf("could not encode block to bytes. Cause: %w", err)
	}
	return encoded, nil
}

func (eb EncodedBlock) Decode() (*types.Block, error) {
	bl := types.Block{}
	err := rlp.DecodeBytes(eb, &bl)
	return &bl, err
}

func (eb EncodedBlock) DecodeBlock() (*types.Block, error) {
	b, err := eb.Decode()
	if err != nil {
		return nil, fmt.Errorf("could not decode block from bytes. Cause: %w", err)
	}
	return b, nil
}

// NewBlock - todo - remove this ASAP
func NewBlock(parent *types.Block, nodeID common.Address, txs []*types.Transaction) *types.Block {
	parentHash := GenesisHash
	height := L1GenesisHeight
	if parent != nil {
		parentHash = parent.Hash()
		height = parent.NumberU64() + 1
	}

	header := types.Header{
		ParentHash:  parentHash,
		UncleHash:   common.Hash{},
		Coinbase:    nodeID,
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  big.NewInt(0),
		Number:      big.NewInt(int64(height)),
		GasLimit:    0,
		GasUsed:     0,
		Time:        0,
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}

	return types.NewBlock(&header, txs, nil, nil, &trie.StackTrie{})
}
