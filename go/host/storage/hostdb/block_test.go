package hostdb

import (
	"database/sql"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestCanStoreAndRetrieveBlock(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	block1 := createBlock(batchNumber)
	block2 := createBlock(batchNumber + 1)

	// verify we get ErrNoRows for a non-existent block
	randomHash := gethcommon.Hash{}
	randomHash.SetBytes(make([]byte, 32)) // 32 bytes for appropriate length
	dbtx, _ := db.NewDBTransaction()
	statements := db.GetSQLStatement()
	_, err := GetBlockId(dbtx.Tx, statements, randomHash)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows for non-existent block, got: %v", err)
	}
	dbtx.Rollback()

	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, statements, &block1)
	if err != nil {
		t.Errorf("could not store block1: %s", err)
	}
	err = dbtx.Write()
	if err != nil {
		t.Errorf("could not commit block1: %s", err)
	}

	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, statements, &block2)
	if err != nil {
		t.Errorf("could not store block2: %s", err)
	}
	err = dbtx.Write()
	if err != nil {
		t.Errorf("could not commit block2: %s", err)
	}

	dbtx, _ = db.NewDBTransaction()
	blockId, err := GetBlockId(dbtx.Tx, statements, block2.Hash())
	if err != nil {
		t.Errorf("stored block but could not retrieve block ID: %s", err)
	}
	if *blockId != 2 {
		t.Errorf("expected block ID 2, got %d", *blockId)
	}
	dbtx.Rollback()
}

func TestAddBlockWithForeignKeyConstraint(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	dbtx, _ := db.NewDBTransaction()
	statements := db.GetSQLStatement()
	metadata := createRollupMetadata(batchNumber - 10)
	rollup := createRollup(batchNumber)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)

	// add block
	err := AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()

	// add rollup referencing block
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup, &common.ExtRollupMetadata{}, &metadata, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	dbtx.Write()

	// this should fail due to the UNIQUE constraint on block_host.hash
	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, statements, block.Header())
	if !strings.Contains(err.Error(), "UNIQUE constraint failed: block_host.hash") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}

	// verify the block still exists
	_, err = GetBlockId(dbtx.Tx, statements, block.Hash())
	if err != nil {
		t.Fatalf("failed to get block id after duplicate insert: %v", err)
	}

	dbtx.Rollback()
}

func createBlock(blockNum int64) types.Header {
	return types.Header{
		Number:     big.NewInt(blockNum),
		Time:       uint64(time.Now().Unix()),
		Difficulty: big.NewInt(1),
		GasLimit:   1000000,
		GasUsed:    0,
		Coinbase:   gethcommon.Address{},
		ParentHash: gethcommon.Hash{},
		Root:       gethcommon.Hash{},
		TxHash:     gethcommon.Hash{},
		ReceiptHash: gethcommon.Hash{},
		Bloom:      types.Bloom{},
		MixDigest:  gethcommon.Hash{},
		Nonce:      types.BlockNonce{},
		BaseFee:    big.NewInt(0),
	}
}

func TestBlockCountInListing(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	statements := db.GetSQLStatement()

	numBlocks := 5
	for i := 0; i < numBlocks; i++ {
		block := createBlock(int64(i + 1))
		dbtx, _ := db.NewDBTransaction()
		err := AddBlock(dbtx.Tx, statements, &block)
		if err != nil {
			t.Errorf("could not store block %d: %s", i+1, err)
		}
		err = dbtx.Write()
		if err != nil {
			t.Errorf("could not commit block %d: %s", i+1, err)
		}
	}

	// verify the total block count is 5
	totalBlocks, err := GetTotalBlockCount(db)
	if err != nil {
		t.Errorf("could not get total block count: %s", err)
	}
	if totalBlocks.Int64() != int64(numBlocks) {
		t.Errorf("expected total block count to be %d, got %d", numBlocks, totalBlocks.Int64())
	}

	pagination := &common.QueryPagination{
		Offset: 0,
		Size:   3,
	}

	blockListing, err := GetBlockListing(db, pagination)
	if err != nil {
		t.Errorf("could not get block listing: %s", err)
	}

	// verify that the total count is correct even when pagination limits the results
	if blockListing.Total != uint64(numBlocks) {
		t.Errorf("expected block listing total to be %d, got %d", numBlocks, blockListing.Total)
	}

	if len(blockListing.BlocksData) != 3 {
		t.Errorf("expected 3 blocks in listing, got %d", len(blockListing.BlocksData))
	}

	pagination2 := &common.QueryPagination{
		Offset: 2,
		Size:   2,
	}

	blockListing2, err := GetBlockListing(db, pagination2)
	if err != nil {
		t.Errorf("could not get block listing with offset: %s", err)
	}

	if blockListing2.Total != uint64(numBlocks) {
		t.Errorf("expected block listing total to be %d, got %d", numBlocks, blockListing2.Total)
	}

	if len(blockListing2.BlocksData) != 2 {
		t.Errorf("expected 2 blocks in listing with offset, got %d", len(blockListing2.BlocksData))
	}

	// Verify that blocks without rollups have zero rollup hashes
	for _, block := range blockListing.BlocksData {
		if block.RollupHash != (common.L2RollupHash{}) {
			t.Errorf("expected blocks without rollups to have zero rollup hash, got %v", block.RollupHash)
		}
	}
}
