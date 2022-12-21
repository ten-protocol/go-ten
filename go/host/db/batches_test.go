package db

import (
	"errors"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/common"
)

func TestCanStoreAndRetrieveBatchHeader(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	batch := common.ExtBatch{
		Header: &header,
	}

	err := db.AddBatchHeader(&batch)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	batchHeader, err := db.GetBatchHeader(header.Hash())
	if err != nil {
		t.Errorf("stored batch header but could not retrieve it. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(header.Number) != 0 {
		t.Errorf("batch header was not stored correctly")
	}
}

func TestUnknownBatchHeaderReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := types.Header{}

	_, err := db.GetBatchHeader(header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch header but was able to retrieve it")
	}
}

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) { //nolint:dupl
	db := NewInMemoryDB(nil)
	headerOne := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	batchOne := common.ExtBatch{
		Header: &headerOne,
	}

	err := db.AddBatchHeader(&batchOne)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	headerTwo := common.BatchHeader{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	batchTwo := common.ExtBatch{
		Header: &headerTwo,
	}

	err = db.AddBatchHeader(&batchTwo)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	batchHeader, err := db.GetHeadBatchHeader()
	if err != nil {
		t.Errorf("stored batch header but could not retrieve it. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(headerTwo.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestLowerNumberBatchDoesNotBecomeBatchHeader(t *testing.T) { //nolint:dupl
	db := NewInMemoryDB(nil)
	headerOne := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	batchOne := common.ExtBatch{
		Header: &headerOne,
	}

	err := db.AddBatchHeader(&batchOne)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	headerTwo := common.BatchHeader{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	batchTwo := common.ExtBatch{
		Header: &headerTwo,
	}

	err = db.AddBatchHeader(&batchTwo)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	batchHeader, err := db.GetHeadBatchHeader()
	if err != nil {
		t.Errorf("stored batch header but could not retrieve it. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(headerOne.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestHeadBatchHeaderIsNotSetInitially(t *testing.T) {
	db := NewInMemoryDB(nil)

	_, err := db.GetHeadBatchHeader()
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head batch was set, but no batchs had been written")
	}
}

func TestCanRetrieveBatchHashByNumber(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	batch := common.ExtBatch{
		Header: &header,
	}

	err := db.AddBatchHeader(&batch)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	batchHash, err := db.GetBatchHash(header.Number)
	if err != nil {
		t.Errorf("stored batch header but could not retrieve its hash by number. Cause: %s", err)
	}
	if *batchHash != header.Hash() {
		t.Errorf("batch hash was not stored correctly against number")
	}
}

func TestUnknownBatchNumberReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := types.Header{}

	_, err := db.GetBatchHash(header.Number)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch hash but was able to retrieve it")
	}
}

func TestCanRetrieveBatchNumberByTxHash(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	batch := common.ExtBatch{
		Header:   &header,
		TxHashes: []gethcommon.Hash{txHash},
	}

	err := db.AddBatchHeader(&batch)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	batchNumber, err := db.GetBatchNumber(txHash)
	if err != nil {
		t.Errorf("stored batch header but could not retrieve its number by transaction hash. Cause: %s", err)
	}
	if batchNumber.Cmp(header.Number) != 0 {
		t.Errorf("batch number was not stored correctly against transaction hash")
	}
}

func TestUnknownBatchTxHashReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB(nil)

	_, err := db.GetBatchNumber(gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveBatchTransactions(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	txHashes := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batch := common.ExtBatch{
		Header:   &header,
		TxHashes: txHashes,
	}

	err := db.AddBatchHeader(&batch)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	batchTxs, err := db.GetBatchTxs(header.Hash())
	if err != nil {
		t.Errorf("stored batch header but could not retrieve its transactions. Cause: %s", err)
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
	db := NewInMemoryDB(nil)

	_, err := db.GetBatchNumber(gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db := NewInMemoryDB(nil)
	headerOne := common.BatchHeader{
		Number: big.NewInt(batchNumber),
	}
	txHashesOne := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := common.ExtBatch{
		Header:   &headerOne,
		TxHashes: txHashesOne,
	}

	err := db.AddBatchHeader(&batchOne)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	headerTwo := common.BatchHeader{
		Number: big.NewInt(batchNumber + 1),
	}
	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := common.ExtBatch{
		Header:   &headerTwo,
		TxHashes: txHashesTwo,
	}

	err = db.AddBatchHeader(&batchTwo)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	totalTxs, err := db.GetTotalTransactions()
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(totalTxs.Int64()) != len(txHashesOne)+len(txHashesTwo) {
		t.Errorf("total number of batch transactions was not stored correctly")
	}
}

// TODO - #718 - Add tests of writing and reading extbatches.
