package playground

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
	"testing"
)

func TestBlockInclusion(t *testing.T) {
	blockchain, _ := NewBlockchain()
	blocks := newChainOfEmptyBlocks(blockchain, blockchain.Genesis())
	panicIfAnyErr(blockchain.InsertChain(blocks))

	assertBlocksIncludedInChain(t, blockchain, blocks)
	assertRandomBlockNotIncludedInChain(t, blockchain)
}

func TestTransactionExecution(t *testing.T) {
	blockchain, db := NewBlockchain()
	key, err := crypto.GenerateKey()
	panicIfErr(err)
	fundAccount(blockchain, key, db)
	blocks := []*types.Block{NewChildBlock(blockchain, blockchain.Genesis(), newTxs(blockchain, key))}
	panicIfAnyErr(blockchain.InsertChain(blocks))

	assertBlocksIncludedInChain(t, blockchain, blocks)
	assertRandomBlockNotIncludedInChain(t, blockchain)

	// TODO - Check transactions are present in blockchain.
}

func newChainOfEmptyBlocks(blockchain *core.BlockChain, firstParent *types.Block) []*types.Block {
	blocks := make([]*types.Block, 5)
	parentBlock := firstParent
	for i := 0; i < 5; i++ {
		block := NewChildBlock(blockchain, parentBlock, []*types.Transaction{})
		blocks[i] = block
		parentBlock = block
	}
	return blocks
}

func fundAccount(blockchain *core.BlockChain, key *ecdsa.PrivateKey, db ethdb.Database) {
	genesisWithPrealloc := core.Genesis{
		Config: core.DefaultGenesisBlock().Config,
		Alloc: map[common.Address]core.GenesisAccount{
			crypto.PubkeyToAddress(key.PublicKey): {Balance: big.NewInt(1000000)},
		},
	}
	// TODO - Can we prealloc at `BlockChain` creation time, rather than setting the genesis block after the fact?
	panicIfErr(blockchain.ResetWithGenesisBlock(genesisWithPrealloc.ToBlock(db)))
}

func newTxs(blockchain *core.BlockChain, key *ecdsa.PrivateKey) []*types.Transaction {
	txs := make([]*types.Transaction, 5)
	for i := 0; i < 5; i++ {
		txData := &types.LegacyTx{
			Nonce: uint64(i),
			Gas:   uint64(21000),
		}
		txs[i] = NewSignedTransaction(blockchain, key, txData)
	}
	return txs
}

func assertBlocksIncludedInChain(t *testing.T, blockchain *core.BlockChain, blocks []*types.Block) {
	for i, block := range blocks {
		if !blockchain.HasBlock(block.Hash(), uint64(i+1)) {
			t.Error("Block was inserted into blockchain, but was not included.")
		}
	}
}

func assertRandomBlockNotIncludedInChain(t *testing.T, blockchain *core.BlockChain) {
	if blockchain.HasBlock(common.BytesToHash([]byte("test_hash")), uint64(0)) {
		t.Errorf("Block was not inserted into blockchain, but was included anyway.")
	}
}
