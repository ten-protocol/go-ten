package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	statedb32 = "statedb32" // the table used for 32 byte keys - 99.9% of the keys are here
	statedb64 = "statedb64" // the table used for larger keys
	getQry    = `select sdb.val from %s sdb where sdb.ky = ? ;`
	// `replace` will perform insert or replace if existing and this syntax works for both sqlite and edgeless db
	putQry       = `replace into %s (ky,  val) values(?,  ?);`
	putQryBatch  = `replace into %s (ky, val) values`
	putQryValues = `(?,?)`
	delQry       = `delete from %s where ky = ? ;`
	// todo - how is the performance of this? probably extraordinarily slow
	searchQry = `select ky, val from %s sdb where substring(sdb.ky, 1, ?) = ? and sdb.ky >= ? order by sdb.ky asc`
)

func getTable(key []byte) string {
	if len(key) <= 32 {
		return statedb32
	}
	return statedb64
}

func Has(ctx context.Context, db *sqlx.DB, key []byte) (bool, error) {
	err := db.QueryRowContext(ctx, fmt.Sprintf(getQry, getTable(key)), key).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Get(ctx context.Context, db *sqlx.DB, key []byte) ([]byte, error) {
	var res []byte

	err := db.QueryRowContext(ctx, fmt.Sprintf(getQry, getTable(key)), key).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

func Put(ctx context.Context, db *sqlx.DB, key []byte, value []byte) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(putQry, getTable(key)), key, value)
	return err
}

func PutKeyValues(ctx context.Context, tx *sqlx.Tx, keys [][]byte, vals [][]byte) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("invalid command. should not happen")
	}

	shortKeys := make([][]byte, 0)
	shortVals := make([][]byte, 0)
	longKeys := make([][]byte, 0)
	longVals := make([][]byte, 0)

	// Split keys and values based on key length
	for i, key := range keys {
		if len(key) <= 32 {
			shortKeys = append(shortKeys, key)
			shortVals = append(shortVals, vals[i])
		} else {
			longKeys = append(longKeys, key)
			longVals = append(longVals, vals[i])
		}
	}

	// Process short keys
	if len(shortKeys) > 0 {
		update := fmt.Sprintf(putQryBatch, statedb32) + repeat(putQryValues, ",", len(shortKeys))
		values := make([]any, 0)
		for i := range shortKeys {
			values = append(values, shortKeys[i], shortVals[i])
		}
		_, err := tx.ExecContext(ctx, update, values...)
		if err != nil {
			return fmt.Errorf("failed to exec short k/v transaction statement. kv=%v, err=%w", values, err)
		}
	}

	// Process long keys
	if len(longKeys) > 0 {
		update := fmt.Sprintf(putQryBatch, statedb64) + repeat(putQryValues, ",", len(longKeys))
		values := make([]any, 0)
		for i := range longKeys {
			values = append(values, longKeys[i], longVals[i])
		}
		_, err := tx.ExecContext(ctx, update, values...)
		if err != nil {
			return fmt.Errorf("failed to exec long k/v transaction statement. kv=%v, err=%w", values, err)
		}
	}

	return nil
}

func Delete(ctx context.Context, db *sqlx.DB, key []byte) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(delQry, getTable(key)), key)
	return err
}

func DeleteKeys(ctx context.Context, db *sqlx.Tx, keys [][]byte) error {
	for _, del := range keys {
		_, err := db.ExecContext(ctx, fmt.Sprintf(delQry, getTable(del)), del)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewIterator(ctx context.Context, db *sqlx.DB, prefix []byte, start []byte) ethdb.Iterator {
	// todo - is this used?
	pr := prefix
	st := append(prefix, start...)

	// iterator clean-up handles closing this rows iterator
	rows, err := db.QueryContext(ctx, fmt.Sprintf(searchQry, getTable(st)), len(pr), pr, st)
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
