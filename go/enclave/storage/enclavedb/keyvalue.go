package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	getQry = `select keyvalue.val from keyvalue where keyvalue.ky = ? and keyvalue.ky_full = ?;`
	// `replace` will perform insert or replace if existing and this syntax works for both sqlite and edgeless db
	putQry       = `replace into keyvalue (ky, ky_full, val) values(?, ?, ?);`
	putQryBatch  = `replace into keyvalue (ky, ky_full, val) values`
	putQryValues = `(?,?,?)`
	delQry       = `delete from keyvalue where keyvalue.ky = ? and keyvalue.ky_full = ?;`
	// todo - how is the performance of this?
	searchQry = `select ky_full, val from keyvalue where substring(keyvalue.ky_full, 1, ?) = ? and keyvalue.ky_full >= ? order by keyvalue.ky_full asc`
)

func Has(ctx context.Context, db *sql.DB, key []byte) (bool, error) {
	err := db.QueryRowContext(ctx, getQry, hash(key), key).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Get(ctx context.Context, db *sql.DB, key []byte) ([]byte, error) {
	var res []byte

	err := db.QueryRowContext(ctx, getQry, hash(key), key).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

func Put(ctx context.Context, db *sql.DB, key []byte, value []byte) error {
	_, err := db.ExecContext(ctx, putQry, hash(key), key, value)
	return err
}

func PutKeyValues(ctx context.Context, tx *sql.Tx, keys [][]byte, vals [][]byte) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("invalid command. should not happen")
	}

	if len(keys) > 0 {
		// write the kv updates as a single update statement for increased efficiency
		update := putQryBatch + repeat(putQryValues, ",", len(keys))

		values := make([]any, 0)
		for i := range keys {
			values = append(values, hash(keys[i]), keys[i], vals[i])
		}
		_, err := tx.ExecContext(ctx, update, values...)
		if err != nil {
			return fmt.Errorf("failed to exec k/v transaction statement. kv=%v, err=%w", values, err)
		}
	}

	return nil
}

func Delete(ctx context.Context, db *sql.DB, key []byte) error {
	_, err := db.ExecContext(ctx, delQry, hash(key), key)
	return err
}

func DeleteKeys(ctx context.Context, db *sql.Tx, keys [][]byte) error {
	for _, del := range keys {
		_, err := db.ExecContext(ctx, delQry, hash(del), del)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewIterator(ctx context.Context, db *sql.DB, prefix []byte, start []byte) ethdb.Iterator {
	pr := prefix
	st := append(prefix, start...)
	// iterator clean-up handles closing this rows iterator
	rows, err := db.QueryContext(ctx, searchQry, len(pr), pr, st)
	if err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}
	if err = rows.Err(); err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}
	return &iterator{
		rows: rows,
	}
}

// hash returns 4 bytes "hash" of the key to be indexed
// truncating is not sufficient because the keys are not random
func hash(key []byte) []byte {
	h := fnv.New32()
	_, err := h.Write(key)
	if err != nil {
		return nil
	}
	return h.Sum([]byte{})
}
