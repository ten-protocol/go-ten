package db

import (
	"errors"
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"
)

// An arbitrary number to put in the header, to check that the header is retrieved correctly from the DB.
const batchNumber = 777

func TestCanStoreAndRetrieveBlockHeader(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := types.Header{
		Number: big.NewInt(batchNumber),
	}
	err := db.AddBlockHeader(&header)
	if err != nil {
		t.Errorf("could not add block header. Cause: %s", err)
	}

	blockHeader, err := db.GetBlockHeader(header.Hash())
	if err != nil {
		t.Errorf("stored block header but could not retrieve it. Cause: %s", err)
	}
	if blockHeader.Number.Cmp(header.Number) != 0 {
		t.Errorf("block header was not stored correctly")
	}
}

func TestUnknownBlockHeaderReturnsNotFound(t *testing.T) {
	db := NewInMemoryDB(nil)
	header := types.Header{}

	_, err := db.GetBlockHeader(header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store block header but was able to retrieve it")
	}
}
