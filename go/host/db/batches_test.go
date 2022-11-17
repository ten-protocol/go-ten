package db

import (
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	db.AddBatchHeader(&headerOne, []gethcommon.Hash{})

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	db.AddBatchHeader(&headerTwo, []gethcommon.Hash{})

	batchHeader, found := db.GetHeadBatchHeader()
	if !found {
		t.Errorf("stored batch header but could not retrieve it")
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
	db.AddBatchHeader(&headerOne, []gethcommon.Hash{})

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	db.AddBatchHeader(&headerTwo, []gethcommon.Hash{})

	batchHeader, found := db.GetHeadBatchHeader()
	if !found {
		t.Errorf("stored batch header but could not retrieve it")
	}
	if batchHeader.Number.Cmp(headerOne.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestHeadBatchHeaderIsNotSetInitially(t *testing.T) {
	db := NewInMemoryDB()

	_, found := db.GetHeadBatchHeader()
	if found {
		t.Errorf("head batch was set, but no batches had been written")
	}
}
