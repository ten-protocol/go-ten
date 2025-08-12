package hostdb

import (
	"errors"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common"
)

// An arbitrary number to put in the header, to check that the header is retrieved correctly from the DB.

func TestCanStoreAndRetrieveBatchHeader(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	batch := CreateBatch(batchNumber, []common.L2TxHash{})
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()
	batchHeader, err := GetBatchHeader(db, batch.Header.Hash())
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch header was not stored correctly")
	}
}

func TestUnknownBatchHeaderReturnsNotFound(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	header := types.Header{}

	_, err := GetBatchHeader(db, header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch header but was able to retrieve it")
	}
}

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) { //nolint:dupl
	db, _ := CreateSQLiteDB(t)
	batchOne := CreateBatch(batchNumber, []common.L2TxHash{})
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTwo := CreateBatch(batchNumber+1, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(dbtx, db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	dbtx.Write()

	batchHeader, err := GetHeadBatchHeader(db)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batchTwo.Header.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestLowerNumberBatchDoesNotBecomeBatchHeader(t *testing.T) { //nolint:dupl
	db, _ := CreateSQLiteDB(t)
	dbtx, _ := db.NewDBTransaction()
	batchOne := CreateBatch(batchNumber, []common.L2TxHash{})
	err := AddBatch(dbtx, db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTwo := CreateBatch(batchNumber-1, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(dbtx, db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	batchHeader, err := GetHeadBatchHeader(db)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batchOne.Header.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestHeadBatchHeaderIsNotSetInitially(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	_, err := GetHeadBatchHeader(db)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head batch was set, but no batchs had been written")
	}
}

func TestCanRetrieveBatchHashByNumber(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	batch := CreateBatch(batchNumber, []common.L2TxHash{})
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	batchHash, err := GetBatchHashByNumber(db, batch.Header.Number)
	if err != nil {
		t.Errorf("stored batch but could not retrieve headers hash by number. Cause: %s", err)
	}
	if *batchHash != batch.Header.Hash() {
		t.Errorf("batch hash was not stored correctly against number")
	}
}

func TestUnknownBatchNumberReturnsNotFound(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	header := types.Header{Number: big.NewInt(10)}

	_, err := GetBatchHashByNumber(db, header.Number)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch hash but was able to retrieve it")
	}
}

func TestCanRetrieveBatchNumberByTxHash(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	batch := CreateBatch(batchNumber, []common.L2TxHash{txHash})
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	extBatch, err := GetBatchByTx(db, txHash)
	if err != nil {
		t.Errorf("stored batch but could not retrieve batch by transaction hash. Cause: %s", err)
	}
	if extBatch.Header.Number.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch number was not stored correctly against transaction hash")
	}
	batchNumber, err := GetBatchNumber(db, txHash)
	if err != nil {
		t.Errorf("stored batch but could not retrieve number by transaction hash. Cause: %s", err)
	}
	if batchNumber.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch number was not stored correctly against transaction hash")
	}
}

func TestUnknownBatchTxHashReturnsNotFound(t *testing.T) {
	db, _ := CreateSQLiteDB(t)

	_, err := GetBatchNumber(db, gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveBatchTransactions(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	txHashes := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batch := CreateBatch(batchNumber, txHashes)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	batchTxs, err := GetBatchTxHashes(db, batch.Header.Hash())
	if err != nil {
		t.Errorf("stored batch but could not retrieve headers transactions. Cause: %s", err)
	}
	if len(batchTxs) != len(txHashes) {
		t.Errorf("batch transactions were not stored correctly")
	}
	for idx, batchTx := range batchTxs {
		if batchTx != txHashes[idx] {
			t.Errorf("batch transactions were not stored correctly")
		}
	}
}

func TestTransactionsForUnknownBatchReturnsNotFound(t *testing.T) {
	db, _ := CreateSQLiteDB(t)

	_, err := GetBatchNumber(db, gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestGetLatestBatch(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := CreateBatch(batchNumber, txHashesOne)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := CreateBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(dbtx, db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	batch, err := GetHeadBatchHeader(db)
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(batch.SequencerOrderNo.Uint64()) != int(batchTwo.SeqNo().Uint64()) {
		t.Errorf("latest batch was not retrieved correctly")
	}
}

func TestGetBatchByHeight(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	batch1 := CreateBatch(batchNumber, []common.L2TxHash{})
	batch2 := CreateBatch(batchNumber+5, []common.L2TxHash{})
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batch1)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	err = AddBatch(dbtx, db, &batch2)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()
	publicBatch, err := GetBatchByHeight(db, batch2.Header.Number)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batch2.Header.Number.Cmp(publicBatch.Header.Number) != 0 {
		t.Errorf("batch header was not stored correctly")
	}
}

func TestGetBatchListing(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := CreateBatch(batchNumber, txHashesOne)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := CreateBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(dbtx, db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesThree := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringFive")), gethcommon.BytesToHash([]byte("magicStringSix"))}
	batchThree := CreateBatch(batchNumber+2, txHashesThree)

	err = AddBatch(dbtx, db, &batchThree)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	// page 1, size 2
	batchListing, err := GetBatchListing(db, &common.QueryPagination{Offset: 1, Size: 2})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be the most recent sequence number
	if big.NewInt(int64(batchListing.Total)).Cmp(big.NewInt(batchNumber+2)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// first element should be the second batch
	if batchListing.BatchesData[0].SequencerOrderNo.Cmp(batchTwo.SeqNo()) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// page 0, size 3
	batchListing1, err := GetBatchListing(db, &common.QueryPagination{Offset: 0, Size: 3})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// first element should be the most recent batch since they're in descending order
	if batchListing1.BatchesData[0].SequencerOrderNo.Cmp(batchThree.SeqNo()) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// page 5, size 1
	batchListing3, err := GetBatchListing(db, &common.QueryPagination{Offset: 5, Size: 1})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be still be the most recent batch sequence number
	if big.NewInt(int64(batchListing3.Total)).Cmp(big.NewInt(batchNumber+2)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}
}

func TestGetBatchTransactions(t *testing.T) {
	db, _ := CreateSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := CreateBatch(batchNumber, txHashesOne)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour")), gethcommon.BytesToHash([]byte("magicStringFive"))}
	batchTwo := CreateBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(dbtx, db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	dbtx.Write()

	txListing1, err := GetBatchTransactions(db, batchOne.Header.Hash(), &common.QueryPagination{Offset: 0, Size: 2})
	if err != nil {
		t.Errorf("stored batch but could not retrieve transactions. Cause: %s", err)
	}

	if txListing1.Total != uint64(len(txHashesOne)) {
		t.Errorf("batch transactions were not retrieved correctly")
	}
	txListing2, err := GetBatchTransactions(db, batchTwo.Header.Hash(), &common.QueryPagination{Offset: 0, Size: 2})
	if err != nil {
		t.Errorf("stored batch but could not retrieve transactions. Cause: %s", err)
	}
	if txListing2.Total != uint64(len(txHashesTwo)) {
		t.Errorf("batch transactions were not retrieved correctly")
	}
}
