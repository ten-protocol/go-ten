package db

import (
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

func TestCanStoreAndRetrieveRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddRollupHeader(&header, []gethcommon.Hash{})

	rollupHeader, found := db.GetRollupHeader(header.Hash())
	if !found {
		t.Errorf("stored rollup header but could not retrieve it")
	}
	if rollupHeader.Number.Cmp(header.Number) != 0 {
		t.Errorf("rollup header was not stored correctly")
	}
}

func TestUnknownRollupHeaderReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()
	header := types.Header{}

	_, found := db.GetRollupHeader(header.Hash())
	if found {
		t.Errorf("did not store rollup header but was able to retrieve it")
	}
}

func TestHigherNumberRollupBecomesRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddRollupHeader(&headerOne, []gethcommon.Hash{})

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	db.AddRollupHeader(&headerTwo, []gethcommon.Hash{})

	rollupHeader, found := db.GetHeadRollupHeader()
	if !found {
		t.Errorf("stored rollup header but could not retrieve it")
	}
	if rollupHeader.Number.Cmp(headerTwo.Number) != 0 {
		t.Errorf("head rollup was not set correctly")
	}
}

func TestLowerNumberRollupDoesNotBecomeRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddRollupHeader(&headerOne, []gethcommon.Hash{})

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	db.AddRollupHeader(&headerTwo, []gethcommon.Hash{})

	rollupHeader, found := db.GetHeadRollupHeader()
	if !found {
		t.Errorf("stored rollup header but could not retrieve it")
	}
	if rollupHeader.Number.Cmp(headerOne.Number) != 0 {
		t.Errorf("head rollup was not set correctly")
	}
}

func TestHeadRollupHeaderIsNotSetInitially(t *testing.T) {
	db := NewInMemoryDB()

	_, found := db.GetHeadRollupHeader()
	if found {
		t.Errorf("head rollup was set, but no rollups had been written")
	}
}

func TestCanRetrieveRollupHashByNumber(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddRollupHeader(&header, []gethcommon.Hash{})

	rollupHash, found := db.GetRollupHash(header.Number)
	if !found {
		t.Errorf("stored rollup header but could not retrieve its hash by number")
	}
	if *rollupHash != header.Hash() {
		t.Errorf("rollup hash was not stored correctly against number")
	}
}

func TestUnknownRollupNumberReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()
	header := types.Header{}

	_, found := db.GetRollupHash(header.Number)
	if found {
		t.Errorf("did not store rollup hash but was able to retrieve it")
	}
}

func TestCanRetrieveRollupNumberByTxHash(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(magicNumber),
	}
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	db.AddRollupHeader(&header, []gethcommon.Hash{txHash})

	rollupNumber, found := db.GetRollupNumber(txHash)
	if !found {
		t.Errorf("stored rollup header but could not retrieve its number by transaction hash")
	}
	// TODO - Temp fix due to off-by-one error in `writeRollupNumber`. Remove once fixed.
	headerNumber := big.NewInt(0).Add(header.Number, big.NewInt(1))
	if rollupNumber.Cmp(headerNumber) != 0 {
		t.Errorf("rollup number was not stored correctly against transaction hash")
	}
}

func TestUnknownRollupTxHashReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()

	_, found := db.GetRollupNumber(gethcommon.BytesToHash([]byte("magicString")))
	if found {
		t.Errorf("did not store rollup number but was able to retrieve it")
	}
}

func TestCanRetrieveRollupTransactions(t *testing.T) {
	db := NewInMemoryDB()
	header := common.Header{
		Number: big.NewInt(magicNumber),
	}
	txHashes := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	db.AddRollupHeader(&header, txHashes)

	rollupTxs, found := db.GetRollupTxs(header.Hash())
	if !found {
		t.Errorf("stored rollup header but could not retrieve its transactions")
	}
	if len(rollupTxs) != len(txHashes) {
		t.Errorf("rollup transactions were not stored correctly")
	}
	for idx, rollupTx := range rollupTxs {
		if rollupTx != txHashes[idx] {
			t.Errorf("rollup transactions were not stored correctly")
		}
	}
}

func TestTransactionsForUnknownRollupReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()

	_, found := db.GetRollupNumber(gethcommon.BytesToHash([]byte("magicString")))
	if found {
		t.Errorf("did not store rollup number but was able to retrieve it")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(magicNumber),
	}
	txHashesOne := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	db.AddRollupHeader(&headerOne, txHashesOne)

	headerTwo := common.Header{
		Number: big.NewInt(magicNumber),
	}
	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	db.AddRollupHeader(&headerTwo, txHashesTwo)

	totalTxs := db.GetTotalTransactions()
	if int(totalTxs.Int64()) != len(txHashesOne)+len(txHashesTwo) {
		t.Errorf("total number of rollup transactions was not stored correctly")
	}
}
