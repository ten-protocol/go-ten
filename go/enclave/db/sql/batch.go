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

type statement struct {
	query string
	args  []any
}

type Batch struct {
	db         *EnclaveDB
	writes     []keyvalue
	statements []statement
	size       int
}

func (b *Batch) ExecuteSQL(query string, args ...any) {
	s := statement{
		query: query,
		args:  args,
	}
	b.statements = append(b.statements, s)
}

// Put inserts the given value into the batch for later committing.
func (b *Batch) Put(key, value []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), common.CopyBytes(value), false})
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *Batch) Delete(key []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), nil, true})
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *Batch) ValueSize() int {
	return b.size
}

// Write executes a batch statement with all the updates
func (b *Batch) Write() error {
	tx, err := b.db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to create batch transaction - %w", err)
	}

	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			_, err = tx.Exec(delQry, keyvalue.key)
			if err != nil {
				return err
			}
		} else {
			_, err = tx.Exec(putQry, keyvalue.key, keyvalue.value)
		}

		if err != nil {
			return fmt.Errorf("failed to exec batch statement. kv=%v, err=%w", keyvalue, err)
		}
	}

	for _, s := range b.statements {
		_, err := tx.Exec(s.query, s.args...)
		if err != nil {
			return fmt.Errorf("failed to exec batch statement. err=%w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit batch of writes - %w", err)
	}
	return nil
}

// Reset resets the batch for reuse.
func (b *Batch) Reset() {
	b.writes = b.writes[:0]
	b.statements = b.statements[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *Batch) Replay(w ethdb.KeyValueWriter) error {
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
