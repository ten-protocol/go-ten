package enclavedb

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	createKVTable = `create table if not exists keyvalue
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    ky      binary(4),
    ky_full varbinary(64),
    val     mediumblob NOT NULL
);`

	key1 = hexutils.HexToBytes("0000000000000000000000000000000000000000000000000000000000000001")
	key2 = hexutils.HexToBytes("0000000000000000000000000000000000000000000000000000000000000002")
	// this key has a different prefix to the others, so we can filter it out of iterator
	diffKey3 = hexutils.HexToBytes("1100000000000000000000000000000000000000000000000000000000000003")
	key4     = hexutils.HexToBytes("0000000000000000000000000000000000000000000000000000000000000004")
)

func TestPutAndGetAndDelHappyPath(t *testing.T) {
	db := createDB(t)
	defer cleanUp(db)

	putData(t, db, key1, []byte("val1"))

	found, err := db.Get(key1)
	failIfError(t, err, "failed to retrieve value")
	assert.Equal(t, found, []byte("val1"))

	err = db.Delete(key1)
	failIfError(t, err, "failed to delete entry")

	_, err = db.Get(key1)
	if err == nil {
		t.Fatal("expected get to fail fetch after deletion")
	}
}

func TestIteratorHappyPath(t *testing.T) {
	db := createDB(t)
	defer cleanUp(db)

	// inserting out of order, we will verify the iterator is ordered by key ascending
	putData(t, db, key4, []byte("val4"))
	putData(t, db, key1, []byte("val1"))
	putData(t, db, key2, []byte("val2"))
	// this doesn't match the key prefix of the others
	putData(t, db, diffKey3, []byte("val3"))

	iter := db.NewIterator(hexutils.HexToBytes("0000"), nil)
	assert.Nil(t, iter.Error())

	assert.True(t, iter.Next())
	assert.Equal(t, key1, iter.Key())
	assert.Equal(t, []byte("val1"), iter.Value())

	assert.True(t, iter.Next())
	assert.Equal(t, key2, iter.Key())
	assert.Equal(t, []byte("val2"), iter.Value())

	assert.True(t, iter.Next())
	assert.Equal(t, key4, iter.Key())
	assert.Equal(t, []byte("val4"), iter.Value())

	//// we expect the diffKey3 entry to have been filtered out by the iterator prefix
	assert.False(t, iter.Next())

	//// create a second iterator that starts from key2 (it should omit key1)
	iterWithFilter := db.NewIterator(hexutils.HexToBytes("000000"), hexutils.HexToBytes("0000000000000000000000000000000000000000000000000000000002"))
	assert.Nil(t, iterWithFilter.Error())

	assert.True(t, iterWithFilter.Next())
	assert.Equal(t, key2, iterWithFilter.Key())
	assert.Equal(t, []byte("val2"), iterWithFilter.Value())

	assert.True(t, iterWithFilter.Next())
	assert.Equal(t, key4, iterWithFilter.Key())
	assert.Equal(t, []byte("val4"), iterWithFilter.Value())
}

func TestBatchUpdateHappyPath(t *testing.T) {
	db := createDB(t)
	defer cleanUp(db)

	batch := db.NewBatch()
	err := batch.Put(key1, []byte("value1"))
	failIfError(t, err, "failed to insert in batch")
	err = batch.Put(key2, []byte("value2"))
	failIfError(t, err, "failed to insert in batch")
	err = batch.Put(diffKey3, []byte("value3"))
	failIfError(t, err, "failed to insert in batch")
	err = batch.Delete(key2)
	failIfError(t, err, "failed to delete from batch")

	err = batch.Write()
	failIfError(t, err, "failed to write batch")

	batch.Reset()

	found, err := db.Get(diffKey3)
	failIfError(t, err, "expected key3 to be in the db")
	assert.Equal(t, found, []byte("value3"))

	_, err = db.Get(key2)
	if err == nil {
		t.Fatal("expected get to fail fetch after deletion")
	}
}

func createDB(t *testing.T) ethdb.Database {
	lite := setupSQLite(t)
	_, err := lite.Exec(createKVTable)
	failIfError(t, err, "Failed to create key-value table in test db")
	s, err := NewEnclaveDB(lite, lite, config.EnclaveConfig{RPCTimeout: time.Second}, testlog.Logger())
	failIfError(t, err, "Failed to create SQLEthDatabase for test")
	return s
}

func cleanUp(db ethdb.Database) {
	_ = db.Close()
}

func putData(t *testing.T, db ethdb.Database, key []byte, val []byte) {
	err := db.Put(key, val)
	failIfError(t, err, "failed to insert into db")
}

func failIfError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatal(msg, err)
	}
}

func setupSQLite(t *testing.T) *sql.DB {
	// create temp sqlite db
	d := t.TempDir()
	f := filepath.Join(d, "test.db")
	db, err := sql.Open("sqlite3", f)
	if err != nil {
		panic(err)
	}
	return db
}
