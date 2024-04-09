package db

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethdb"
	ethldb "github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

// ObscuroLevelDB is a very thin wrapper around a level DB database for compatibility with our internal interfaces
// In particular, it overrides the Get method to return the TEN ErrNotFound
type ObscuroLevelDB struct {
	db *ethldb.Database
}

func (o *ObscuroLevelDB) NewBatchWithSize(int) ethdb.Batch {
	// TODO implement me
	panic("implement me")
}

func (o *ObscuroLevelDB) NewSnapshot() (ethdb.Snapshot, error) {
	// TODO implement me
	panic("implement me")
}

func (o *ObscuroLevelDB) Has(key []byte) (bool, error) {
	return o.db.Has(key)
}

// Get is overridden here to return our internal NotFound error
func (o *ObscuroLevelDB) Get(key []byte) ([]byte, error) {
	d, err := o.db.Get(key)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return d, nil
}

func (o *ObscuroLevelDB) Put(key []byte, value []byte) error {
	return o.db.Put(key, value)
}

func (o *ObscuroLevelDB) Delete(key []byte) error {
	return o.db.Delete(key)
}

func (o *ObscuroLevelDB) NewBatch() ethdb.Batch {
	return o.db.NewBatch()
}

func (o *ObscuroLevelDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return o.db.NewIterator(prefix, start)
}

func (o *ObscuroLevelDB) Stat(property string) (string, error) {
	return o.db.Stat(property)
}

func (o *ObscuroLevelDB) Compact(start []byte, limit []byte) error {
	return o.db.Compact(start, limit)
}

func (o *ObscuroLevelDB) Close() error {
	return o.db.Close()
}
