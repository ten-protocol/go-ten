package orm

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	obscurosql "github.com/obscuronet/go-obscuro/go/enclave/db/sql"
)

const cfgInsert = "insert into config values (?,?)"
const cfgUpdate = "update config set val=? where key=?"
const cfgSelect = "select val from config where key=?"

const attInsert = "insert into attestation_key values (?,?)"
const attSelect = "select key from attestation_key where party=?"

func WriteConfigToBatch(dbtx *obscurosql.Batch, key string, value []byte) {
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
	rows, err := db.Query(query, v)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errutil.ErrNotFound
	}
	var r sql.RawBytes
	err = rows.Scan(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
