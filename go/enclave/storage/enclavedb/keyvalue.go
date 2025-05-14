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
	getQry = `select keyvalue.val from keyvalue where keyvalue.ky = ? ;`
	// `replace` will perform insert or replace if existing and this syntax works for both sqlite and edgeless db
	putQry       = `replace into keyvalue (ky,  val) values(?,  ?);`
	putQryBatch  = `replace into keyvalue (ky, val) values`
	putQryValues = `(?,?)`
	delQry       = `delete from keyvalue where keyvalue.ky = ? ;`
	// todo - how is the performance of this?
	searchQry = `select ky, val from keyvalue where substring(keyvalue.ky, 1, ?) = ? and keyvalue.ky >= ? order by keyvalue.ky asc`
)

func Has(ctx context.Context, db *sqlx.DB, key []byte) (bool, error) {
	err := db.QueryRowContext(ctx, getQry, key).Scan()
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

	err := db.QueryRowContext(ctx, getQry, key).Scan(&res)
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
	_, err := db.ExecContext(ctx, putQry, key, value)
	return err
}

func PutKeyValues(ctx context.Context, tx *sqlx.Tx, keys [][]byte, vals [][]byte) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("invalid command. should not happen")
	}

	if len(keys) > 0 {
		// write the kv updates as a single update statement for increased efficiency
		update := putQryBatch + repeat(putQryValues, ",", len(keys))

		values := make([]any, 0)
		for i := range keys {
			values = append(values, keys[i], vals[i])
		}
		_, err := tx.ExecContext(ctx, update, values...)
		if err != nil {
			return fmt.Errorf("failed to exec k/v transaction statement. kv=%v, err=%w", values, err)
		}
	}

	return nil
}

func Delete(ctx context.Context, db *sqlx.DB, key []byte) error {
	_, err := db.ExecContext(ctx, delQry, key)
	return err
}

func DeleteKeys(ctx context.Context, db *sqlx.Tx, keys [][]byte) error {
	for _, del := range keys {
		_, err := db.ExecContext(ctx, delQry, del)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewIterator(ctx context.Context, db *sqlx.DB, prefix []byte, start []byte) ethdb.Iterator {
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
