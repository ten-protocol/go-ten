package enclavedb

import (
	"context"
	"database/sql"
	"encoding/json"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
)

// NetworkUpgrade represents a network upgrade feature
type NetworkUpgrade struct {
	ID                uint64                 `db:"id" json:"id"`
	FeatureName       string                 `db:"feature_name" json:"feature_name"`
	FeatureData       []byte                 `db:"feature_data" json:"-"`
	FeatureDataMap    map[string]interface{} `db:"-" json:"feature_data"`
	BlockHash         gethcommon.Hash        `db:"block_hash" json:"block_hash"`
	BlockHeightFinal  *uint64                `db:"block_height_final" json:"block_height_final"`
	BlockHeightActive *uint64                `db:"block_height_active" json:"block_height_active"`
	CreatedAt         sql.NullTime           `db:"created_at" json:"created_at"`
}

// MarshalFeatureData converts FeatureDataMap to FeatureData bytes for database storage
func (nu *NetworkUpgrade) MarshalFeatureData() error {
	if nu.FeatureDataMap != nil {
		data, err := json.Marshal(nu.FeatureDataMap)
		if err != nil {
			return err
		}
		nu.FeatureData = data
	}
	return nil
}

// UnmarshalFeatureData converts FeatureData bytes to FeatureDataMap for easier usage
func (nu *NetworkUpgrade) UnmarshalFeatureData() error {
	if len(nu.FeatureData) > 0 {
		return json.Unmarshal(nu.FeatureData, &nu.FeatureDataMap)
	}
	return nil
}

// StoreNetworkUpgrade stores a new network upgrade in the database
func StoreNetworkUpgrade(ctx context.Context, db *sqlx.Tx, upgrade *NetworkUpgrade) error {
	// Ensure feature data is marshaled
	if err := upgrade.MarshalFeatureData(); err != nil {
		return err
	}

	query := `INSERT INTO network_upgrade 
				(feature_name, feature_data, block_hash, block_height_final, block_height_active) 
			  VALUES (?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query,
		upgrade.FeatureName,
		upgrade.FeatureData,
		upgrade.BlockHash.Bytes(),
		upgrade.BlockHeightFinal,
		upgrade.BlockHeightActive,
	)

	return err
}

// GetNetworkUpgrades retrieves all network upgrades from the database
func GetNetworkUpgrades(ctx context.Context, db *sqlx.DB) ([]*NetworkUpgrade, error) {
	query := `SELECT id, feature_name, feature_data, block_hash, block_height_final, block_height_active, created_at 
			  FROM network_upgrade 
			  ORDER BY created_at ASC`

	var upgrades []*NetworkUpgrade
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var upgrade NetworkUpgrade
		var blockHashBytes []byte

		err := rows.Scan(
			&upgrade.ID,
			&upgrade.FeatureName,
			&upgrade.FeatureData,
			&blockHashBytes,
			&upgrade.BlockHeightFinal,
			&upgrade.BlockHeightActive,
			&upgrade.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert block_hash bytes to gethcommon.Hash
		upgrade.BlockHash = gethcommon.BytesToHash(blockHashBytes)

		// Unmarshal feature data
		if err := upgrade.UnmarshalFeatureData(); err != nil {
			return nil, err
		}

		upgrades = append(upgrades, &upgrade)
	}

	return upgrades, rows.Err()
}

// GetNetworkUpgradesByFeature retrieves network upgrades for a specific feature
func GetNetworkUpgradesByFeature(ctx context.Context, db *sqlx.DB, featureName string) ([]*NetworkUpgrade, error) {
	query := `SELECT id, feature_name, feature_data, block_hash, block_height_final, block_height_active, created_at 
			  FROM network_upgrade 
			  WHERE feature_name = ?
			  ORDER BY created_at ASC`

	var upgrades []*NetworkUpgrade
	rows, err := db.QueryxContext(ctx, query, featureName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var upgrade NetworkUpgrade
		var blockHashBytes []byte

		err := rows.Scan(
			&upgrade.ID,
			&upgrade.FeatureName,
			&upgrade.FeatureData,
			&blockHashBytes,
			&upgrade.BlockHeightFinal,
			&upgrade.BlockHeightActive,
			&upgrade.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert block_hash bytes to gethcommon.Hash
		upgrade.BlockHash = gethcommon.BytesToHash(blockHashBytes)

		// Unmarshal feature data
		if err := upgrade.UnmarshalFeatureData(); err != nil {
			return nil, err
		}

		upgrades = append(upgrades, &upgrade)
	}

	return upgrades, rows.Err()
}

func GetActivatedNetworkUpgrades(ctx context.Context, db *sqlx.DB, height uint64) ([]*NetworkUpgrade, error) {
	query := `SELECT id, feature_name, feature_data, block_hash, block_height_final, block_height_active, created_at
			  FROM network_upgrade
			  WHERE block_height_active IS NOT NULL AND block_height_active <= ?
			  ORDER BY created_at ASC`

	var upgrades []*NetworkUpgrade
	rows, err := db.QueryxContext(ctx, query, height)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var upgrade NetworkUpgrade
		var blockHashBytes []byte

		err := rows.Scan(
			&upgrade.ID,
			&upgrade.FeatureName,
			&upgrade.FeatureData,
			&blockHashBytes,
			&upgrade.BlockHeightFinal,
			&upgrade.BlockHeightActive,
			&upgrade.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		upgrade.BlockHash = gethcommon.BytesToHash(blockHashBytes)
		if err := upgrade.UnmarshalFeatureData(); err != nil {
			return nil, err
		}
		upgrades = append(upgrades, &upgrade)
	}

	return upgrades, rows.Err()
}

// DeleteNetworkUpgrade removes a network upgrade from the database
func DeleteNetworkUpgrade(ctx context.Context, db *sqlx.Tx, id uint64) error {
	query := `DELETE FROM network_upgrade WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// UpdateNetworkUpgradeHeights updates the finalization and activation heights for a network upgrade
func UpdateNetworkUpgradeHeights(ctx context.Context, db *sqlx.Tx, id uint64, finalHeight, activeHeight *uint64) error {
	query := `UPDATE network_upgrade 
			  SET block_height_final = ?, block_height_active = ? 
			  WHERE id = ?`

	_, err := db.ExecContext(ctx, query, finalHeight, activeHeight, id)
	return err
}
