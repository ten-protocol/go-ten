package database

import (
	"database/sql"
	"errors"

	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	cfgInsert = "insert into config values (?,?)"
	cfgSelect = "select val from config where ky=?"
)

func WriteConfigToTx(dbtx *sql.Tx, key string, value any) (sql.Result, error) {
	return dbtx.Exec(cfgInsert, key, value)
}

func WriteConfig(db *sql.DB, key string, value []byte) (sql.Result, error) {
	return db.Exec(cfgInsert, key, value)
}

func FetchConfig(db *sql.DB, key string) ([]byte, error) {
	return readSingleRow(db, cfgSelect, key)
}

func readSingleRow(db *sql.DB, query string, v any) ([]byte, error) {
	var res []byte

	err := db.QueryRow(query, v).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}
