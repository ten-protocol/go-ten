package db

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
)

const magicNumber = 777

func TestCanStoreAndRetrieveBlockHeader(t *testing.T) {
	db := NewInMemoryDB()
	header := types.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddBlockHeader(&header)

	blockHeader, found := db.GetBlockHeader(header.Hash())
	if !found {
		t.Errorf("stored block header but could not retrieve it")
	}
	if blockHeader.Number.Cmp(header.Number) != 0 {
		t.Errorf("block header was not stored correctly")
	}
}

func TestUnknownBlockHeaderReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB()
	header := types.Header{}

	_, found := db.GetBlockHeader(header.Hash())
	if found {
		t.Errorf("did not store block header but was able to retrieve it")
	}
}

func TestHigherNumberBlockBecomesBlockHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := types.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddBlockHeader(&headerOne)

	headerTwo := types.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
	}
	db.AddBlockHeader(&headerTwo)

	blockHeader, found := db.GetHeadBlockHeader()
	if !found {
		t.Errorf("stored block header but could not retrieve it")
	}
	if blockHeader.Number.Cmp(headerTwo.Number) != 0 {
		t.Errorf("head block was not set correctly")
	}
}

func TestLowerNumberBlockDoesNotBecomeBlockHeader(t *testing.T) {
	db := NewInMemoryDB()
	headerOne := types.Header{
		Number: big.NewInt(magicNumber),
	}
	db.AddBlockHeader(&headerOne)

	headerTwo := types.Header{
		// We give the second header a higher number, making it the head.
		Number: big.NewInt(0).Sub(headerOne.Number, big.NewInt(1)),
	}
	db.AddBlockHeader(&headerTwo)

	blockHeader, found := db.GetHeadBlockHeader()
	if !found {
		t.Errorf("stored block header but could not retrieve it")
	}
	if blockHeader.Number.Cmp(headerOne.Number) != 0 {
		t.Errorf("head block was not set correctly")
	}
}

func TestHeadBlockHeaderIsNotSetInitially(t *testing.T) {
	db := NewInMemoryDB()

	_, found := db.GetHeadBlockHeader()
	if found {
		t.Errorf("head block was set, but no blocks had been written")
	}
}
