package orm

import (
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	obscurosql "github.com/obscuronet/go-obscuro/go/enclave/db/sql"
)

const (
	cfgInsert = "insert into config values (?,?)"
	cfgUpdate = "update config set val=? where ky=?"
	cfgSelect = "select val from config where ky=?"
)

const (
	attInsert = "insert into attestation_key values (?,?)"
	attSelect = "select ky from attestation_key where party=?"
)

func WriteConfigToBatch(dbtx *obscurosql.Batch, key string, value any) {
	dbtx.ExecuteSQL(cfgInsert, key, value)
}

func WriteConfig(db *sql.DB, key string, value []byte) (sql.Result, error) {
	return db.Exec(cfgInsert, key, value)
}

func UpdateConfigToBatch(dbtx *obscurosql.Batch, key string, value []byte) {
	dbtx.ExecuteSQL(cfgUpdate, key, value)
}

func UpdateConfig(db *sql.DB, key string, value []byte) (sql.Result, error) {
	return db.Exec(cfgUpdate, key, value)
}

func FetchConfig(db *sql.DB, key string) ([]byte, error) {
	return readSingleRow(db, cfgSelect, key)
}

func WriteAttKey(db *sql.DB, party common.Address, key []byte) (sql.Result, error) {
	return db.Exec(attInsert, party.Bytes(), key)
}

func FetchAttKey(db *sql.DB, party common.Address) ([]byte, error) {
	return readSingleRow(db, attSelect, party.Bytes())
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
