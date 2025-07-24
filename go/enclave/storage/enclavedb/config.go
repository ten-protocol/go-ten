package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	cfgExists = "select exists(select 1 from config where ky = ?)"
	cfgInsert = "insert into config values (?,?)"
	cfgUpdate = "update config set val = ? where ky = ?"
	cfgSelect = "select val from config where ky=?"
)

const (
	attInsert           = "insert into attestation (enclave_id, pub_key, node_type)  values (?,?,?)"
	attSelect           = "select pub_key, node_type from attestation where enclave_id=?"
	attExists           = "select exists(select 1 from attestation where enclave_id = ?)"
	attUpdateKey        = "update attestation set pub_key = ? where enclave_id = ?"
	attUpdate           = "update attestation set node_type=? where enclave_id=?"
	attSelectSequencers = "select enclave_id from attestation where node_type = ?"
)

func InsertOrUpdateConfig(ctx context.Context, dbtx *sqlx.Tx, key string, value any) error {
	var exists bool
	// check if it exists then insert or update - this keeps it agnostic to the type of sql database
	err := dbtx.GetContext(ctx, &exists, cfgExists, key)
	if err != nil {
		return fmt.Errorf("failed to check existence of config key %q: %w", key, err)
	}

	if exists {
		_, err = dbtx.ExecContext(ctx, cfgUpdate, value, key)
		if err != nil {
			return fmt.Errorf("failed to update config key %q: %w", key, err)
		}
	} else {
		_, err = dbtx.ExecContext(ctx, cfgInsert, key, value)
		if err != nil {
			return fmt.Errorf("failed to insert config key %q: %w", key, err)
		}
	}

	return nil
}

func WriteConfig(ctx context.Context, db *sqlx.Tx, key string, value []byte) (sql.Result, error) {
	return db.ExecContext(ctx, cfgInsert, key, value)
}

func FetchConfig(ctx context.Context, db *sqlx.DB, key string) ([]byte, error) {
	return readSingleRow(ctx, db, cfgSelect, key)
}

func WriteAttestation(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID, key []byte, nodeType common.NodeType) (sql.Result, error) {
	return db.ExecContext(ctx, attInsert, enclaveId.Bytes(), key, nodeType)
}

func UpdateAttestationKey(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID, key []byte) (sql.Result, error) {
	return db.ExecContext(ctx, attUpdateKey, key, enclaveId.Bytes())
}

func UpdateAttestationType(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID, nodeType common.NodeType) (sql.Result, error) {
	return db.ExecContext(ctx, attUpdate, nodeType, enclaveId.Bytes())
}

func FetchAttestation(ctx context.Context, db *sqlx.DB, enclaveId common.EnclaveID) ([]byte, common.NodeType, error) {
	var pubKey []byte
	var nodeType common.NodeType

	err := db.QueryRowContext(ctx, attSelect, enclaveId.Bytes()).Scan(&pubKey, &nodeType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, 0, errutil.ErrNotFound
		}
		return nil, 0, err
	}
	return pubKey, nodeType, nil
}

// AttestationExists checks if an attestation exists for the given enclave ID
func AttestationExists(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, attExists, enclaveId.Bytes()).Scan(&exists)
	return exists, err
}

func readSingleRow(ctx context.Context, db *sqlx.DB, query string, v any) ([]byte, error) {
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

// FetchSequencerEnclaveIDs returns all enclave IDs that are registered as sequencers
func FetchSequencerEnclaveIDs(ctx context.Context, db *sqlx.DB) ([]common.EnclaveID, error) {
	rows, err := db.QueryContext(ctx, attSelectSequencers, common.Sequencer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enclaveIDs []common.EnclaveID
	for rows.Next() {
		var idBytes []byte
		if err := rows.Scan(&idBytes); err != nil {
			return nil, err
		}
		enclaveID := common.EnclaveID{}
		enclaveID.SetBytes(idBytes)
		enclaveIDs = append(enclaveIDs, enclaveID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return enclaveIDs, nil
}
