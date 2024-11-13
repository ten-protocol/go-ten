package hostdb

import (
	"database/sql"
	"errors"
	"math/big"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestCanStoreAndRetrieveBlock(t *testing.T) {
	db, _ := createSQLiteDB(t)
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

func createBlock(blockNum int64) types.Header {
	return types.Header{
		Number: big.NewInt(blockNum),
		Time:   uint64(time.Now().Unix()),
	}
}
