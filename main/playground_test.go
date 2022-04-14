package playground

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)

//func TestBlockInclusion(t *testing.T) {
//	blockchain, _ := createBlockchain()
//
//	blocks := make([]*types.Block, 5)
//	parentBlock := blockchain.Genesis()
//	for i := 0; i < 5; i++ {
//		block := newChildBlock(parentBlock, nil)
//		blocks[i] = block
//		parentBlock = block
//	}
//
//	_, err := blockchain.InsertChain(blocks)
//	panicIfErr(err)
//
//	// Check all inserted blocks are included.
//	for i, block := range blocks {
//		if !blockchain.HasBlock(block.Hash(), uint64(i+1)) {
//			t.Error("Block was inserted into blockchain, but was not included.")
//		}
//	}
//
//	// Check a random non-inserted block isn't included.
//	if blockchain.HasBlock(common.BytesToHash([]byte("test_hash")), uint64(0)) {
//		t.Errorf("Block was not inserted into blockchain, but was included anyway.")
//	}
//}

func TestTransactionExecution(t *testing.T) {
	blockchain, db := createBlockchain()

	key, err := crypto.GenerateKey()
	panicIfErr(err)
	address := crypto.PubkeyToAddress(key.PublicKey)
	account := core.GenesisAccount{Balance: big.NewInt(1000000)}

	genesisWithPrealloc := core.Genesis{
		Config: core.DefaultGenesisBlock().Config,
		Alloc: map[common.Address]core.GenesisAccount{
			address: account,
		},
	}

	err = blockchain.ResetWithGenesisBlock(genesisWithPrealloc.ToBlock(db)) // todo - joel - is there a better way to prealloc?
	panicIfErr(err)

	txData := &types.LegacyTx{
		Nonce: uint64(0),
		Gas:   uint64(21000),
	}
	tx := newSignedTransaction(blockchain, key, txData)

	parentBlock := blockchain.Genesis()

	// I have to create the block once with no receipts, in order to produce the receipts, in order to add the receipts to the block.
	blockForReceipts := newChildBlock(parentBlock, []*types.Transaction{tx}, nil)
	stateDb, err := blockchain.State()
	panicIfErr(err)
	receipts, _, _, err := blockchain.Processor().Process(blockForReceipts, stateDb, vm.Config{})
	panicIfErr(err)

	//blockchain.State()
	block := newChildBlock(parentBlock, []*types.Transaction{tx}, receipts)

	_, err = blockchain.InsertChain([]*types.Block{block})
	panicIfErr(err)
}
