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
const batchNumber = 777

func TestCanStoreAndRetrieveBatchHeader(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	batch, err := getBatch(batchNumber, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHeader, err := GetBatchHeader(db, batch.Header.Hash())
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch header was not stored correctly")
	}
}

func TestUnknownBatchHeaderReturnsNotFound(t *testing.T) {
	db, err := createSQLiteDB(t)
	header := types.Header{}

	_, err = GetBatchHeader(db, header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch header but was able to retrieve it")
	}
}

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) { //nolint:dupl
	db, err := createSQLiteDB(t)
	batchOne, err := getBatch(batchNumber, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTwo, err := getBatch(batchNumber+1, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHeader, err := GetHeadBatchHeader(db)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batchTwo.Header.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestLowerNumberBatchDoesNotBecomeBatchHeader(t *testing.T) { //nolint:dupl
	db, err := createSQLiteDB(t)
	batchOne, err := getBatch(batchNumber, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTwo, err := getBatch(batchNumber-1, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHeader, err := GetHeadBatchHeader(db)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batchOne.Header.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestHeadBatchHeaderIsNotSetInitially(t *testing.T) {
	//FIXME
	db, err := createSQLiteDB(t)

	_, err = GetHeadBatchHeader(db)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head batch was set, but no batchs had been written")
	}
}

func TestCanRetrieveBatchHashByNumber(t *testing.T) {
	//FIXME Implement me
	db, err := createSQLiteDB(t)
	batch, err := getBatch(batchNumber, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHash, err := GetBatchHash(db, batch.Header.Number)
	if err != nil {
		t.Errorf("stored batch but could not retrieve headers hash by number. Cause: %s", err)
	}
	if *batchHash != batch.Header.Hash() {
		t.Errorf("batch hash was not stored correctly against number")
	}
}

func TestUnknownBatchNumberReturnsNotFound(t *testing.T) {
	db, err := createSQLiteDB(t)
	header := types.Header{}

	_, err = GetBatchHash(db, header.Number)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch hash but was able to retrieve it")
	}
}

func TestCanRetrieveBatchNumberByTxHash(t *testing.T) {
	db, err := createSQLiteDB(t)
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	batch, err := getBatch(batchNumber, []common.L2TxHash{txHash})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchNumber, err := GetBatchNumber(db, txHash)
	if err != nil {
		t.Errorf("stored batch but could not retrieve headers number by transaction hash. Cause: %s", err)
	}
	if batchNumber.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch number was not stored correctly against transaction hash")
	}
}

func TestUnknownBatchTxHashReturnsNotFound(t *testing.T) {
	db, err := createSQLiteDB(t)

	_, err = GetBatchNumber(db, gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveBatchTransactions(t *testing.T) {
	db, err := createSQLiteDB(t)
	txHashes := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batch, err := getBatch(batchNumber, txHashes)
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTxs, err := GetBatchTxs(db, batch.Header.Hash())
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
	db, err := createSQLiteDB(t)

	_, err = GetBatchNumber(db, gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db, err := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne, err := getBatch(batchNumber, txHashesOne)

	err = AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo, err := getBatch(batchNumber, txHashesTwo)

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	totalTxs, err := GetTotalTransactions(db)
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(totalTxs.Int64()) != len(txHashesOne)+len(txHashesTwo) {
		t.Errorf("total number of batch transactions was not stored correctly")
	}
}

func getBatch(batchNum int64, txHashes []common.L2BatchHash) (common.ExtBatch, error) {
	header := common.BatchHeader{
		SequencerOrderNo: big.NewInt(batchNum),
		Number:           big.NewInt(batchNum),
	}
	batch := common.ExtBatch{
		Header:   &header,
		TxHashes: txHashes,
	}

	return batch, nil
}

// todo (#718) - add tests of writing and reading extbatches.
