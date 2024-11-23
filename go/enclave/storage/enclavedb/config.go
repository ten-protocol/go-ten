package enclavedb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	cfgInsert = "insert into config values (?,?)"
	cfgSelect = "select val from config where ky=?"
)

const (
	attInsert = "insert into attestation (enclave_id, pub_key, is_sequencer)  values (?,?,?)"
	attSelect = "select pub_key, is_sequencer from attestation where enclave_id=?"
	attUpdate = "update attestation set is_sequencer=? where enclave_id=?"
)

func WriteConfigToTx(ctx context.Context, dbtx *sql.Tx, key string, value any) (sql.Result, error) {
	return dbtx.ExecContext(ctx, cfgInsert, key, value)
}

func WriteConfig(ctx context.Context, db *sql.Tx, key string, value []byte) (sql.Result, error) {
	return db.ExecContext(ctx, cfgInsert, key, value)
}

func FetchConfig(ctx context.Context, db *sql.DB, key string) ([]byte, error) {
	return readSingleRow(ctx, db, cfgSelect, key)
}

func WriteAttestation(ctx context.Context, db *sql.Tx, enclaveId common.EnclaveID, key []byte, isSequencer bool) (sql.Result, error) {
	return db.ExecContext(ctx, attInsert, enclaveId.Bytes(), key, isSequencer)
}

func UpdateAttestation(ctx context.Context, db *sql.Tx, enclaveId common.EnclaveID, isSequencer bool) (sql.Result, error) {
	return db.ExecContext(ctx, attUpdate, enclaveId.Bytes(), isSequencer)
}

func FetchAttestation(ctx context.Context, db *sql.DB, enclaveId common.EnclaveID) ([]byte, bool, error) {
	var pubKey []byte
	var isSequencer bool

	err := db.QueryRowContext(ctx, attSelect, enclaveId.Bytes()).Scan(&pubKey, &isSequencer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, false, errutil.ErrNotFound
		}
		return nil, false, err
	}
	return pubKey, isSequencer, nil
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
