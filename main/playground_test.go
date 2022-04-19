package playground

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestBlockInclusion(t *testing.T) {
	blockchain, _ := NewBlockchain()
	txs := make([][]*types.Transaction, 5)
	for i := 0; i < 5; i++ {
		txs[i] = []*types.Transaction{}
	}
	blocks := NewChainOfBlocks(blockchain, blockchain.Genesis(), txs)
	panicIfAnyErr(blockchain.InsertChain(blocks))

	assertBlocksIncludedInChain(t, blockchain, blocks)
	assertRandomBlockNotIncludedInChain(t, blockchain)
}

func TestTransactionInclusion(t *testing.T) {
	blockchain, db := NewBlockchain()
	key, err := crypto.GenerateKey()
	panicIfErr(err)
	PrefundKeys(blockchain, []*ecdsa.PrivateKey{key}, db)

	// When handcrafting transactions, we have to create and insert each block in turn. We cannot prepare a series of
	// blocks, then insert them all at once. This is because we use `BlockChain.Processor().Process` when creating a
	// block to generate the tx receipts for us. If we don't update the blockchain after each block creation, `Process`
	// will expect each block to use tx nonce starting from the same initial value. But this nonce reuse will then be
	// rejected when we attempt to insert the blocks into the chain.
	blocks := make([]*types.Block, 3)
	txsPerBlock := make([][]*types.Transaction, 3)
	for i := 0; i < 3; i++ {
		txs := NewTxs(blockchain, key, 5)
		block := NewChildBlock(blockchain, blockchain.Genesis(), txs)
		panicIfAnyErr(blockchain.InsertChain([]*types.Block{block}))
		blocks[i] = block
		txsPerBlock[i] = txs
	}

	assertBlocksIncludedInChain(t, blockchain, blocks)
	assertTxsIncludedInChain(t, blockchain, blocks, txsPerBlock)
	assertRandomBlockNotIncludedInChain(t, blockchain)
}

func assertBlocksIncludedInChain(t *testing.T, blockchain *core.BlockChain, blocks []*types.Block) {
	for i, block := range blocks {
		if !blockchain.HasBlock(block.Hash(), uint64(i+1)) {
			t.Error("Block was inserted into blockchain, but was not included.")
		}
	}
}

func assertTxsIncludedInChain(t *testing.T, blockchain *core.BlockChain, blocks []*types.Block, txsPerBlock [][]*types.Transaction) {
	for i, block := range blocks {
		retrievedBlock := blockchain.GetBlockByHash(block.Hash())

		for _, tx := range txsPerBlock[i] {
			if retrievedBlock.Transaction(tx.Hash()) == nil {
				t.Error("Transactions were inserted into blockchain, but were not included.")
			}
		}
	}
}

func assertRandomBlockNotIncludedInChain(t *testing.T, blockchain *core.BlockChain) {
	if blockchain.HasBlock(common.BytesToHash([]byte("test_hash")), uint64(0)) {
		t.Errorf("Block was not inserted into blockchain, but was included anyway.")
	}
}
