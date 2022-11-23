package db

import (
	"errors"
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

func TestHigherNumberRollupBecomesRollupHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := common.Header{
		Number: big.NewInt(rollupNumber),
	}
	err := db.AddRollupHeader(&headerOne)
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddRollupHeader(&headerTwo)
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
	err := db.AddRollupHeader(&headerOne)
	if err != nil {
		t.Errorf("could not store rollup header. Cause: %s", err)
	}

	headerTwo := common.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	err = db.AddRollupHeader(&headerTwo)
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
	err := db.AddRollupHeader(&header)
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
