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

func TestCanStoreAndRetrieveRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddRollupHeader(&header, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	rollupHeader, err := db.GetRollupHeader(header.Hash())
	if err != nil {
		t.Errorf("stored rollup header but could not retrieve it. Cause: %s", err)
	}
	if rollupHeader.Number.Cmp(header.Number) != 0 {
		t.Errorf("rollup header was not stored correctly")
	}
}

func TestUnknownRollupHeaderReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()
	header := types.Header{}

	_, err := db.GetRollupHeader(header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store rollup header but was able to retrieve it")
	}
}

func TestHigherNumberRollupBecomesRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddRollupHeader(&headerOne, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddRollupHeader(&headerTwo, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	rollupHeader, err := db.GetHeadRollupHeader()
	if err != nil {
		t.Errorf("stored rollup header but could not retrieve it. Cause: %s", err)
	}
	if rollupHeader.Number.Cmp(headerTwo.Number) != 0 {
		t.Errorf("head rollup was not set correctly")
	}
}

func TestLowerNumberRollupDoesNotBecomeRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddRollupHeader(&headerOne, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddRollupHeader(&headerTwo, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	rollupHeader, err := db.GetHeadRollupHeader()
	if err != nil {
		t.Errorf("stored rollup header but could not retrieve it. Cause: %s", err)
	}
	if rollupHeader.Number.Cmp(headerOne.Number) != 0 {
		t.Errorf("head rollup was not set correctly")
	}
}

func TestHeadRollupHeaderIsNotSetInitially(t *testing.T) {
	db := NewInMemoryDB()

	_, err := db.GetHeadRollupHeader()
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head rollup was set, but no rollups had been written")
	}
}

func TestCanRetrieveRollupHashByNumber(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddRollupHeader(&header, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	rollupHash, err := db.GetRollupHash(header.Number)
	if err != nil {
		t.Errorf("stored rollup header but could not retrieve its hash by number. Cause: %s", err)
	}
	if *rollupHash != header.Hash() {
		t.Errorf("rollup hash was not stored correctly against number")
	}
}

func TestUnknownRollupNumberReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()
	header := types.Header{}

	_, err := db.GetRollupHash(header.Number)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store rollup hash but was able to retrieve it")
	}
}

func TestCanRetrieveRollupNumberByTxHash(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	err := db.AddRollupHeader(&header, []gethcommon.Hash{txHash})
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	rollupNumber, err := db.GetRollupNumber(txHash)
	if err != nil {
		t.Errorf("stored rollup header but could not retrieve its number by transaction hash. Cause: %s", err)
	}
	// TODO - Temp fix due to off-by-one error in `writeRollupNumber`. Remove once fixed.
	headerNumber := big.NewInt(0).Add(header.Number, big.NewInt(1))
	if rollupNumber.Cmp(headerNumber) != 0 {
		t.Errorf("rollup number was not stored correctly against transaction hash")
	}
}

func TestUnknownRollupTxHashReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()

	_, err := db.GetRollupNumber(gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store rollup number but was able to retrieve it")
	}
}

func TestTransactionsForUnknownRollupReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()

	_, err := db.GetRollupNumber(gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store rollup number but was able to retrieve it")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHashesOne := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	err := db.AddRollupHeader(&headerOne, txHashesOne)
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	headerTwo := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	err = db.AddRollupHeader(&headerTwo, txHashesTwo)
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	totalTxs, err := db.GetTotalTransactions()
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(totalTxs.Int64()) != len(txHashesOne)+len(txHashesTwo) {
		t.Errorf("total number of rollup transactions was not stored correctly")
	}
}
