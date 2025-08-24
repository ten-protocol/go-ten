package enclavedb

import (
	"context"
	"database/sql"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
)

const (
	StatusPending   = "pending"
	StatusFinalized = "finalized"
)

// NetworkUpgrade represents a network upgrade event stored in the database
type NetworkUpgrade struct {
	ID                uint64
	FeatureName       string
	FeatureData       []byte
	AppliedAtL1Height uint64
	AppliedAtL1Hash   gethcommon.Hash
	TxHash            gethcommon.Hash
	Status            string // "pending" or "finalized"
	FinalizedAtHeight *uint64
	FinalizedAtHash   *gethcommon.Hash
	CreatedAt         *string
}

// WriteNetworkUpgrade stores a new pending network upgrade in the database
func WriteNetworkUpgrade(ctx context.Context, dbtx *sqlx.Tx, featureName string, featureData []byte, l1Height uint64, l1Hash gethcommon.Hash, txHash gethcommon.Hash) error {
	insertSQL := "INSERT INTO network_upgrades (feature_name, feature_data, applied_at_l1_height, applied_at_l1_hash, tx_hash, status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := dbtx.ExecContext(ctx, insertSQL, featureName, featureData, l1Height, l1Hash.Hex(), txHash.Hex(), StatusPending)
	if err != nil {
		return fmt.Errorf("failed to insert network upgrade: %w", err)
	}
	return nil
}

// FinalizeNetworkUpgrade marks a pending network upgrade as finalized
func FinalizeNetworkUpgrade(ctx context.Context, dbtx *sqlx.Tx, txHash gethcommon.Hash, finalizedAtHeight uint64, finalizedAtHash gethcommon.Hash) error {
	updateSQL := "UPDATE network_upgrades SET status = ?, finalized_at_height = ?, finalized_at_hash = ? WHERE tx_hash = ? AND status = ?"
	result, err := dbtx.ExecContext(ctx, updateSQL, StatusFinalized, finalizedAtHeight, finalizedAtHash.Hex(), txHash.Hex(), StatusPending)
	if err != nil {
		return fmt.Errorf("failed to finalize network upgrade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no pending network upgrade found with tx hash %s", txHash.Hex())
	}

	return nil
}

// ReadNetworkUpgradesByStatus returns all network upgrades with the specified status
func ReadNetworkUpgradesByStatus(ctx context.Context, db *sqlx.DB, status string) ([]NetworkUpgrade, error) {
	selectSQL := `SELECT id, feature_name, feature_data, applied_at_l1_height, applied_at_l1_hash, tx_hash, status, 
	                     finalized_at_height, finalized_at_hash, created_at 
	              FROM network_upgrades 
	              WHERE status = ? 
	              ORDER BY applied_at_l1_height ASC, id ASC`

	rows, err := db.QueryContext(ctx, selectSQL, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query network upgrades: %w", err)
	}
	defer rows.Close()

	var upgrades []NetworkUpgrade
	for rows.Next() {
		var upgrade NetworkUpgrade
		var l1Hash, txHash string
		var finalizedAtHeight sql.NullInt64
		var finalizedAtHash sql.NullString
		var createdAt sql.NullString

		err := rows.Scan(
			&upgrade.ID,
			&upgrade.FeatureName,
			&upgrade.FeatureData,
			&upgrade.AppliedAtL1Height,
			&l1Hash,
			&txHash,
			&upgrade.Status,
			&finalizedAtHeight,
			&finalizedAtHash,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan network upgrade row: %w", err)
		}

		// Convert hex strings to hashes
		upgrade.AppliedAtL1Hash = gethcommon.HexToHash(l1Hash)
		upgrade.TxHash = gethcommon.HexToHash(txHash)

		// Handle nullable fields
		if finalizedAtHeight.Valid {
			height := uint64(finalizedAtHeight.Int64)
			upgrade.FinalizedAtHeight = &height
		}
		if finalizedAtHash.Valid {
			hash := gethcommon.HexToHash(finalizedAtHash.String)
			upgrade.FinalizedAtHash = &hash
		}
		if createdAt.Valid {
			upgrade.CreatedAt = &createdAt.String
		}

		upgrades = append(upgrades, upgrade)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating network upgrade rows: %w", err)
	}

	return upgrades, nil
}

// ReadPendingNetworkUpgrades returns all pending network upgrades
func ReadPendingNetworkUpgrades(ctx context.Context, db *sqlx.DB) ([]NetworkUpgrade, error) {
	return ReadNetworkUpgradesByStatus(ctx, db, StatusPending)
}

// ReadFinalizedNetworkUpgrades returns all finalized network upgrades
func ReadFinalizedNetworkUpgrades(ctx context.Context, db *sqlx.DB) ([]NetworkUpgrade, error) {
	return ReadNetworkUpgradesByStatus(ctx, db, StatusFinalized)
}
