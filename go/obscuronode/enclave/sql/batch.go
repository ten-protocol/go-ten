package sql

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

// ---- Implement the geth Batch interface, re-using ideas and types from geth's memorydb.go ----

// keyvalue is a key-value tuple that can be flagged with a deletion field to allow creating database write batches.
type keyvalue struct {
	key    []byte
	value  []byte
	delete bool
}

type sqlBatch struct {
	db     *sqlEthDatabase
	writes []keyvalue
	size   int
}

// Put inserts the given value into the batch for later committing.
func (b *sqlBatch) Put(key, value []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), common.CopyBytes(value), false})
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *sqlBatch) Delete(key []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), nil, true})
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *sqlBatch) ValueSize() int {
	return b.size
}

// Write executes a batch statement with all the updates
func (b *sqlBatch) Write() error {
	tx, err := b.db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to prepare batch transaction - %w", err)
	}
	pIns, err := tx.Prepare(`insert or replace into kv values(?, ?);`)
	defer func() { _ = pIns.Close() }()
	if err != nil {
		return fmt.Errorf("failed to create batch prepared stmt - %w", err)
	}
	pDel, err := tx.Prepare(`delete from kv where kv.key = ?;`)
	defer func() { _ = pDel.Close() }()
	if err != nil {
		return fmt.Errorf("failed to create batch prepared stmt - %w", err)
	}

	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			_, err = pDel.Exec(keyvalue.key)
			if err != nil {
				return err
			}
		} else {
			_, err = pIns.Exec(keyvalue.key, keyvalue.value)
		}

		if err != nil {
			return fmt.Errorf("failed to exec batch statement. kv=%v, err=%w", keyvalue, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit batch of writes - %w", err)
	}
	return nil
}

// Reset resets the batch for reuse.
func (b *sqlBatch) Reset() {
	b.writes = b.writes[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *sqlBatch) Replay(w ethdb.KeyValueWriter) error {
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
