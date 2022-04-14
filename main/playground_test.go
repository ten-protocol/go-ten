package playground

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
)

func TestBlockInclusion(t *testing.T) {
	blockchain, db := createBlockchain()

	blocks := make([]*types.Block, 5)
	parentBlock := core.DefaultGenesisBlock().ToBlock(db)
	for i := 0; i < 5; i++ {
		block := newChildBlock(parentBlock, nil)
		blocks[i] = block
		parentBlock = block
	}

	_, err := blockchain.InsertChain(blocks)
	panicIfErr(err)

	// Check all inserted blocks are included.
	for i, block := range blocks {
		if !blockchain.HasBlock(block.Hash(), uint64(i)) {
			t.Error("Block was inserted into blockchain, but was not included.")
		}
	}

	// Check a random non-inserted block isn't included.
	if blockchain.HasBlock(common.BytesToHash([]byte("test_hash")), uint64(0)) {
		t.Errorf("Block was not inserted into blockchain, but was included anyway.")
	}
}

func TestTransactionExecution(t *testing.T) {
	blockchain, db := createBlockchain()

	// todo - joel - need to prefund some gas
	txData := &types.LegacyTx{
		Nonce: uint64(0),
		Gas:   uint64(21000),
	}
	tx := newSignedTransaction(blockchain, txData)

	parentBlock := core.DefaultGenesisBlock().ToBlock(db)
	block := newChildBlock(parentBlock, []*types.Transaction{tx})

	_, err := blockchain.InsertChain([]*types.Block{block})
	panicIfErr(err)
}
