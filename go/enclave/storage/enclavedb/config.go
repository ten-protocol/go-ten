package enclavedb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	attInsert           = "insert into attestation (enclave_id, pub_key, node_type)  values (?,?,?)"
	attSelect           = "select pub_key, node_type from attestation where enclave_id=?"
	attExists           = "select exists(select 1 from attestation where enclave_id = ?)"
	attUpdateKey        = "update attestation set pub_key = ? where enclave_id = ?"
	attUpdate           = "update attestation set node_type=? where enclave_id=?"
	attSelectSequencers = "select enclave_id from attestation where node_type = ?"
)

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
