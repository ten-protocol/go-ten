package enclavedb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	attInsert         = "insert into attestation (enclave_id, pub_key, report, node_type)  values (?,?,?,?)"
	attSelect         = "select pub_key, node_type from attestation where enclave_id=?"
	attExists         = "select exists(select 1 from attestation where enclave_id = ?)"
	attUpdateKey      = "update attestation set pub_key = ? where enclave_id = ?"
	attUpdate         = "update attestation set node_type=? where enclave_id=?"
	attSelectEnclaves = "select enclave_id from attestation where node_type = ?"
	attSelectReports  = "select enclave_id, pub_key, report from attestation where node_type = ?"
)

func WriteAttestation(ctx context.Context, db *sqlx.Tx, attestation common.AttestationReport, key []byte, nodeType common.NodeType) (sql.Result, error) {
	return db.ExecContext(ctx, attInsert, attestation.EnclaveID.Bytes(), key, attestation.Report, nodeType)
}

func UpdateAttestationKey(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID, key []byte) (sql.Result, error) {
	return db.ExecContext(ctx, attUpdateKey, key, enclaveId.Bytes())
}

func UpdateAttestationType(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID, nodeType common.NodeType) (sql.Result, error) {
	return db.ExecContext(ctx, attUpdate, nodeType, enclaveId.Bytes())
}

// FetchPublicKey returns the public key of the enclave from the attestation table
func FetchPublicKey(ctx context.Context, db *sqlx.DB, enclaveId common.EnclaveID) ([]byte, common.NodeType, error) {
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

// FetchSequencerAttestations returns all attestation reports for enclaves that are sequencers
func FetchSequencerAttestations(ctx context.Context, db *sqlx.DB) ([]common.AttestationReport, error) {
	rows, err := db.QueryContext(ctx, attSelectReports, common.Sequencer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attestations []common.AttestationReport
	for rows.Next() {
		var idBytes []byte
		var pubKey []byte
		var report []byte
		if err := rows.Scan(&idBytes, &pubKey, &report); err != nil {
			return nil, err
		}
		addr := gethcommon.BytesToAddress(idBytes)
		attestations = append(attestations, common.AttestationReport{
			Report:    report,
			PubKey:    pubKey,
			EnclaveID: addr,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attestations, nil
}

// AttestationExists checks if an attestation exists for the given enclave ID
func AttestationExists(ctx context.Context, db *sqlx.Tx, enclaveId common.EnclaveID) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, attExists, enclaveId.Bytes()).Scan(&exists)
	return exists, err
}

// FetchSequencerEnclaveIDs returns all enclave IDs that are registered as sequencers
func FetchSequencerEnclaveIDs(ctx context.Context, db *sqlx.DB) ([]common.EnclaveID, error) {
	rows, err := db.QueryContext(ctx, attSelectEnclaves, common.Sequencer)
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

// FetchValidatorEnclaveIDs returns all enclave IDs that are registered as validators (node_type = 0)
func FetchValidatorEnclaveIDs(ctx context.Context, db *sqlx.DB) ([]common.EnclaveID, error) {
	rows, err := db.QueryContext(ctx, attSelectEnclaves, common.Validator)
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
