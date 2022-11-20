package db

import (
	"errors"
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddBatchHeader(&headerOne, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not add batch header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddBatchHeader(&headerTwo, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not add batch header. Cause: %s", err)
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
		t.Errorf("could not add batch header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddBatchHeader(&headerTwo, []gethcommon.Hash{})
	if err != nil {
		t.Errorf("could not add batch header. Cause: %s", err)
	}

	batchHeader, err := db.GetHeadBatchHeader()
	if err != nil {
		t.Errorf("stored head batch header but could not retrieve it. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(headerOne.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestHeadBatchHeaderIsNotSetInitially(t *testing.T) {
	db := NewInMemoryDB()

	_, err := db.GetHeadBatchHeader()
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head batch was set, but no batches had been written")
	}
}
