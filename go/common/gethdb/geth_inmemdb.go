package gethdb

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

type InMemGethDB struct {
	*memorydb.Database
}

func NewMemDB() ethdb.KeyValueStore {
	return &InMemGethDB{memorydb.New()}
}

// Get retrieves the given key if it's present in the key-value store.
// this method is adapted to ensure the correct error is always returned
func (db *InMemGethDB) Get(key []byte) ([]byte, error) {
	value, err := db.Database.Get(key)
	if err != nil {
		// memoryDB will return not found error instead of a custom error type
		if err.Error() == "not found" {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return value, nil
}
