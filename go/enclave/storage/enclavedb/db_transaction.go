package enclavedb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

// todo - adjust this value
var deadline = 5 * time.Second

// ---- Implement the geth Batch interface, re-using ideas and types from geth's memorydb.go ----

// keyvalue is a key-value tuple that can be flagged with a deletion field to allow creating database write batches.
type keyvalue struct {
	key    []byte
	value  []byte
	delete bool
}

type statement struct {
	query string
	args  []any
}

type dbTransaction struct {
	db         EnclaveDB
	writes     []keyvalue
	statements []statement
	size       int
}

func (b *dbTransaction) GetDB() *sql.DB {
	return b.db.GetSQLDB()
}

func (b *dbTransaction) ExecuteSQL(query string, args ...any) {
	s := statement{
		query: query,
		args:  args,
	}
	b.statements = append(b.statements, s)
}

// Put inserts the given value into the batch for later committing.
func (b *dbTransaction) Put(key, value []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), common.CopyBytes(value), false})
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *dbTransaction) Delete(key []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), nil, true})
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *dbTransaction) ValueSize() int {
	return b.size
}

// Write executes a batch statement with all the updates
func (b *dbTransaction) Write() error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), deadline)
	defer cancelCtx()
	return b.WriteCtx(ctx)
}

func (b *dbTransaction) WriteCtx(ctx context.Context) error {
	tx, err := b.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to create batch transaction - %w", err)
	}

	var deletes [][]byte
	var updateKeys [][]byte
	var updateValues [][]byte

	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			deletes = append(deletes, keyvalue.key)
		} else {
			updateKeys = append(updateKeys, keyvalue.key)
			updateValues = append(updateValues, keyvalue.value)
		}
	}

	err = PutKeyValues(ctx, tx, updateKeys, updateValues)
	if err != nil {
		return fmt.Errorf("failed to put key/value. Cause %w", err)
	}

	err = DeleteKeys(ctx, tx, deletes)
	if err != nil {
		return fmt.Errorf("failed to delete keys. Cause %w", err)
	}

	for _, s := range b.statements {
		_, err := tx.Exec(s.query, s.args...)
		if err != nil {
			return fmt.Errorf("failed to exec db statement `%s` (%v). Cause: %w", s.query, s.args, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit batch of writes. Cause: %w", err)
	}
	return nil
}

// Reset resets the batch for reuse.
func (b *dbTransaction) Reset() {
	b.writes = b.writes[:0]
	b.statements = b.statements[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *dbTransaction) Replay(w ethdb.KeyValueWriter) error {
	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			if err := w.Delete(keyvalue.key); err != nil {
				return err
			}
			continue
		}
		if err := w.Put(keyvalue.key, keyvalue.value); err != nil {
			return err
		}
	}
	return nil
}
