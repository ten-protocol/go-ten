package db

import (
	"errors"
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

func TestCanStoreAndRetrieveBatchHeader(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddBatchHeader(&header, []gethcommon.Hash{})
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
	db := NewInMemoryDB()
	header := types.Header{}

	_, err := db.GetBatchHeader(header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch header but was able to retrieve it")
	}
}

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddBatchHeader(&headerOne, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddBatchHeader(&headerTwo, []gethcommon.Hash{})
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

func TestLowerNumberBatchDoesNotBecomeBatchHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddBatchHeader(&headerOne, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddBatchHeader(&headerTwo, []gethcommon.Hash{})
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
	db := NewInMemoryDB()

	_, err := db.GetHeadBatchHeader()
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head batch was set, but no batchs had been written")
	}
}

func TestCanRetrieveBatchHashByNumber(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddBatchHeader(&header, []gethcommon.Hash{})
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
	db := NewInMemoryDB()
	header := types.Header{}

	_, err := db.GetBatchHash(header.Number)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch hash but was able to retrieve it")
	}
}

func TestCanRetrieveBatchNumberByTxHash(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	err := db.AddBatchHeader(&header, []gethcommon.Hash{txHash})
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
	db := NewInMemoryDB()

	_, err := db.GetBatchNumber(gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveBatchTransactions(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHashes := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	err := db.AddBatchHeader(&header, txHashes)
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
	db := NewInMemoryDB()

	_, err := db.GetBatchNumber(gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHashesOne := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	err := db.AddBatchHeader(&headerOne, txHashesOne)
	if err != nil {
		t.Errorf("could not store batch header. Cause: %s", err)
	}

	headerTwo := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	err = db.AddBatchHeader(&headerTwo, txHashesTwo)
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
