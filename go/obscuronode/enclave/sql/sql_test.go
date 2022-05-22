package sql

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestPutAndGetAndDelHappyPath(t *testing.T) {
	db := createDB(t)
	defer cleanUp(db)

	putData(t, db, "key1", "val1")

	found, err := db.Get([]byte("key1"))
	failIfError(t, err, "failed to retrieve value")
	assert.Equal(t, found, []byte("val1"))

	err = db.Delete([]byte("key1"))
	failIfError(t, err, "failed to delete entry")

	_, err = db.Get([]byte("key1"))
	if err == nil {
		t.Fatal("expected get to fail fetch after deletion")
	}
}

func TestIteratorHappyPath(t *testing.T) {
	db := createDB(t)
	defer cleanUp(db)

	putData(t, db, "key1", "val1")
	putData(t, db, "key2", "val2")
	putData(t, db, "diffKey3", "val3")

	iter := db.NewIterator([]byte("key"), nil)
	assert.Nil(t, iter.Error())

	assert.True(t, iter.Next())
	assert.Equal(t, []byte("key1"), iter.Key())
	assert.Equal(t, []byte("val1"), iter.Value())

	assert.True(t, iter.Next())
	assert.Equal(t, []byte("key2"), iter.Key())
	assert.Equal(t, []byte("val2"), iter.Value())

	// we expect the other entry to be filtered by the iterator prefix
	assert.False(t, iter.Next())
}

func TestBatchUpdateHappyPath(t *testing.T) {
	db := createDB(t)
	defer cleanUp(db)

	batch := db.NewBatch()
	err := batch.Put([]byte("key1"), []byte("value1"))
	failIfError(t, err, "failed to insert in batch")
	err = batch.Put([]byte("key2"), []byte("value2"))
	failIfError(t, err, "failed to insert in batch")
	err = batch.Put([]byte("key3"), []byte("value3"))
	failIfError(t, err, "failed to insert in batch")
	err = batch.Delete([]byte("key2"))
	failIfError(t, err, "failed to delete from batch")

	err = batch.Write()
	failIfError(t, err, "failed to write batch")

	batch.Reset()

	found, err := db.Get([]byte("key3"))
	failIfError(t, err, "expected key3 to be in the db")
	assert.Equal(t, found, []byte("value3"))

	_, err = db.Get([]byte("key2"))
	if err == nil {
		t.Fatal("expected get to fail fetch after deletion")
	}
}

func createDB(t *testing.T) ethdb.Database {
	// todo: is it valid to use sqlite for testing when we'll be using a mysql-based db?
	lite := setupSQLite(t)
	s, err := CreateSQLEthDatabase(lite)
	if err != nil {
		panic(err)
	}
	return s
}

func cleanUp(db ethdb.Database) {
	_ = db.Close()
}

func putData(t *testing.T, db ethdb.Database, key string, val string) {
	err := db.Put([]byte(key), []byte(val))
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
