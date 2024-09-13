package ethereummock

import (
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

var MockGenesisBlock, _ = NewBlock(nil, common.HexToAddress("0x0"), []*types.Transaction{})

func NewBlock(parent *types.Block, nodeID common.Address, txs []*types.Transaction) (*types.Block, []kzg4844.Blob) {
	var parentHash common.Hash
	var height uint64
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
	var blobs []kzg4844.Blob

	for _, tx := range txs {
		if tx.BlobHashes() != nil {
			blobs = append(blobs, tx.BlobTxSidecar().Blobs...)
		}
	}

	return types.NewBlock(&header, &types.Body{Transactions: txs}, nil, trie.NewStackTrie(nil)), blobs
}
