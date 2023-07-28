package enclavedb

import (
	"database/sql"
	"log"
)

// ---- iterator mostly ported from geth's memorydb.go ----

// iterator can walk over the (potentially partial) keyspace of a memory key
// value store. Internally it is a wrapper of the sql result iterator
type iterator struct {
	rows    *sql.Rows
	currKey string
	currVal []byte
	err     error
}

// Next calls next on the sql Rows iterator and if there was a next then it sets the curr values for reading
func (it *iterator) Next() bool {
	next := it.err == nil && it.rows.Next()
	if !next {
		it.currKey = ""
		it.currVal = []byte{}
		return false
	}

	err := it.rows.Scan(&it.currKey, &it.currVal)
	if err != nil {
		log.Printf("failed to scan row in sql iterator - %s", err)
		it.err = err
	}
	return true
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error. A memory iterator cannot encounter errors.
func (it *iterator) Error() error {
	return it.err
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (it *iterator) Key() []byte {
	return []byte(it.currKey)
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (it *iterator) Value() []byte {
	return it.currVal
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (it *iterator) Release() {
	if it.rows != nil {
		_ = it.rows.Close()
		it.currKey = ""
		it.currVal = []byte{}
	}
}
