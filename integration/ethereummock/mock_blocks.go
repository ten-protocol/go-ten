package ethereummock

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

var MockGenesisBlock = NewBlock(nil, common.HexToAddress("0x0"), []*types.Transaction{}, 0)

func NewBlock(parent *types.Block, nodeID common.Address, txs []*types.Transaction, blockTime uint64) *types.Block {
	var parentHash common.Hash
	var height uint64
	if parent != nil {
		parentHash = parent.Hash()
		height = parent.NumberU64() + 1
	}

	blobGasUsed := uint64(0)
	excessBlobGas := uint64(0)
	header := types.Header{
		ParentHash:       parentHash,
		UncleHash:        common.Hash{},
		Coinbase:         nodeID,
		Root:             common.Hash{},
		TxHash:           common.Hash{},
		ReceiptHash:      common.Hash{},
		Bloom:            types.Bloom{},
		Difficulty:       big.NewInt(0),
		Number:           big.NewInt(int64(height)),
		GasLimit:         100,
		GasUsed:          100,
		Time:             blockTime, // Set the block time here
		Extra:            make([]byte, 0),
		MixDigest:        common.Hash{},
		Nonce:            types.BlockNonce{},
		BaseFee:          big.NewInt(1),
		WithdrawalsHash:  &common.Hash{},
		BlobGasUsed:      &blobGasUsed,
		ExcessBlobGas:    &excessBlobGas,
		ParentBeaconRoot: &common.Hash{},
	}
	return types.NewBlock(&header, &types.Body{Transactions: txs}, nil, trie.NewStackTrie(nil))
}
