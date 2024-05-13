package enclavedb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
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

func WriteConfigToTx(ctx context.Context, dbtx *sql.Tx, key string, value any) (sql.Result, error) {
	return dbtx.Exec(cfgInsert, key, value)
}

func WriteConfig(ctx context.Context, db *sql.Tx, key string, value []byte) (sql.Result, error) {
	return db.ExecContext(ctx, cfgInsert, key, value)
}

func FetchConfig(ctx context.Context, db *sql.DB, key string) ([]byte, error) {
	return readSingleRow(ctx, db, cfgSelect, key)
}

func WriteAttKey(ctx context.Context, db *sql.Tx, party common.Address, key []byte) (sql.Result, error) {
	return db.ExecContext(ctx, attInsert, party.Bytes(), key)
}

func FetchAttKey(ctx context.Context, db *sql.DB, party common.Address) ([]byte, error) {
	return readSingleRow(ctx, db, attSelect, party.Bytes())
}

func readSingleRow(ctx context.Context, db *sql.DB, query string, v any) ([]byte, error) {
	var res []byte

	err := db.QueryRowContext(ctx, query, v).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}
