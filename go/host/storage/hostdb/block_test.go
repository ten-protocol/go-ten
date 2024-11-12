package hostdb

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
	"time"
)

func TestCanStoreAndRetrieveBlock(t *testing.T) {
	db, _ := createSQLiteDB(t)
	block1 := createBlock(batchNumber)
	block2 := createBlock(batchNumber + 1)
	dbtx, _ := db.NewDBTransaction()
	err := AddBlock(dbtx, db.GetSQLStatement(), &block1)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	err = AddBlock(dbtx, db.GetSQLStatement(), &block2)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	blockId, err := GetBlockId(db, block2.Hash())
	if err != nil {
		t.Errorf("stored block but could not retrieve block ID. Cause: %s", err)
	}
	if *blockId != 2 {
		t.Errorf("block was not stored correctly")
	}
}

func createBlock(blockNum int64) types.Header {
	return types.Header{
		Number: big.NewInt(blockNum),
		Time:   uint64(time.Now().Unix()),
	}
}
