package enclavedb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
)

const (
	getQry = `select keyvalue.val from keyvalue where keyvalue.ky = ?;`
	// `replace` will perform insert or replace if existing and this syntax works for both sqlite and edgeless db
	putQry       = `replace into keyvalue values(?, ?);`
	putQryBatch  = `replace into keyvalue values`
	putQryValues = `(?,?)`
	delQry       = `delete from keyvalue where keyvalue.ky = ?;`
	searchQry    = `select * from keyvalue where substring(keyvalue.ky, 1, ?) = ? and keyvalue.ky >= ? order by keyvalue.ky asc`
)

func Has(db *sql.DB, key []byte) (bool, error) {
	err := db.QueryRow(getQry, key).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Get(db *sql.DB, key []byte) ([]byte, error) {
	var res []byte

	err := db.QueryRow(getQry, key).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

func Put(db *sql.DB, key []byte, value []byte) error {
	_, err := db.Exec(putQry, key, value)
	return err
}

func PutKeyValues(tx *sql.Tx, keys [][]byte, vals [][]byte) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("invalid command. should not happen")
	}

	if len(keys) > 0 {
		// write the kv updates as a single update statement for increased efficiency
		update := putQryBatch + strings.Repeat(putQryValues+",", len(keys))
		values := make([]any, 0)
		for i := range keys {
			values = append(values, keys[i], vals[i])
		}
		_, err := tx.Exec(update[0:len(update)-1], values...)
		if err != nil {
			return fmt.Errorf("failed to exec batch statement. kv=%v, err=%w", values, err)
		}
	}

	return nil
}

func Delete(db *sql.DB, key []byte) error {
	_, err := db.Exec(delQry, key)
	return err
}

func DeleteKeys(db *sql.Tx, keys [][]byte) error {
	for _, del := range keys {
		_, err := db.Exec(delQry, del)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewIterator(db *sql.DB, prefix []byte, start []byte) ethdb.Iterator {
	pr := prefix
	st := append(prefix, start...)
	// iterator clean-up handles closing this rows iterator
	rows, err := db.Query(searchQry, len(pr), pr, st)
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
