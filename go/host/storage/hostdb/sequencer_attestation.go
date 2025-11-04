package hostdb

import (
	"database/sql"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const (
	insertSequencerAttestation = "INSERT INTO sequencer_attestation_host (enclave_id, is_active) VALUES (?, ?)"
	updateSequencerActive      = "UPDATE sequencer_attestation_host SET is_active = ? WHERE enclave_id = ?"
	selectAllActiveSequencers  = "SELECT enclave_id FROM sequencer_attestation_host WHERE is_active = 1"
)

// AddSequencerAttestation stores a sequencer enclave ID in the host DB
func AddSequencerAttestation(dbtx *sql.Tx, db *sqlx.DB, enclaveID gethcommon.Address, isActive bool) error {
	reboundInsert := db.Rebind(insertSequencerAttestation)
	_, err := dbtx.Exec(reboundInsert, enclaveID.Bytes(), isActive)
	if err != nil {
		if IsRowExistsError(err) {
			return errutil.ErrAlreadyExists
		}
		return fmt.Errorf("could not insert sequencer attestation. Cause: %w", err)
	}
	return nil
}

// UpdateSequencerStatus updates the active status of a sequencer
func UpdateSequencerStatus(dbtx *sql.Tx, db *sqlx.DB, enclaveID gethcommon.Address, isActive bool) error {
	reboundUpdate := db.Rebind(updateSequencerActive)
	result, err := dbtx.Exec(reboundUpdate, isActive, enclaveID.Bytes())
	if err != nil {
		return fmt.Errorf("could not update sequencer status. Cause: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errutil.ErrNotFound
	}

	return nil
}

// GetAllActiveSequencers retrieves all active sequencer enclave IDs
func GetAllActiveSequencers(db *sqlx.DB) ([]gethcommon.Address, error) {
	reboundSelect := db.Rebind(selectAllActiveSequencers)
	rows, err := db.Query(reboundSelect)
	if err != nil {
		return nil, fmt.Errorf("could not query active sequencers. Cause: %w", err)
	}
	defer rows.Close()

	var sequencers []gethcommon.Address
	for rows.Next() {
		var enclaveIDBytes []byte
		if err := rows.Scan(&enclaveIDBytes); err != nil {
			return nil, fmt.Errorf("could not scan sequencer row. Cause: %w", err)
		}
		sequencers = append(sequencers, gethcommon.BytesToAddress(enclaveIDBytes))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sequencer rows. Cause: %w", err)
	}

	return sequencers, nil
}
