package hostdb

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestCanStoreAndRetrieveBlock(t *testing.T) {
	db, _ := createSQLiteDB(t)
	block1 := createBlock(batchNumber)
	block2 := createBlock(batchNumber + 1)

	dbtx, _ := db.NewDBTransaction()
	err := AddBlock(dbtx, db.GetSQLStatement(), &block1)
	if err != nil {
		t.Errorf("could not store block1: %s", err)
	}
	err = dbtx.Write()
	if err != nil {
		t.Errorf("could not commit block1: %s", err)
	}

	// second block, new tx
	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx, db.GetSQLStatement(), &block2)
	if err != nil {
		t.Errorf("could not store block2: %s", err)
	}
	err = dbtx.Write()
	if err != nil {
		t.Errorf("could not commit block2: %s", err)
	}

	blockId, err := GetBlockId(db, block2.Hash())
	if err != nil {
		t.Errorf("stored block but could not retrieve block ID: %s", err)
	}
	// sqlite indexes from 1 with auto increment
	if *blockId != 2 {
		t.Errorf("expected block ID 2, got %d", *blockId)
	}
}

func createBlock(blockNum int64) types.Header {
	return types.Header{
		Number: big.NewInt(blockNum),
		Time:   uint64(time.Now().Unix()),
	}
}
