package storage

import (
	"fmt"

	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/host/storage/init/postgres"
)

const HOST = "HOST_"

// CreateDBFromConfig creates an appropriate ethdb.Database instance based on your config
func CreateDBFromConfig(cfg *hostconfig.HostConfig, logger gethlog.Logger) (hostdb.HostDB, error) {
	dbName := HOST + cfg.ID
	var db *sqlx.DB
	var err error
	if err = validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		// create sqlite db
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		db, err = sqlite.CreateTemporarySQLiteHostDB(dbName, "mode=memory&cache=shared&_foreign_keys=on", logger)
		if err != nil {
			return nil, fmt.Errorf("could not create in memory sqlite DB: %w", err)
		}
	} else {
		// create postgres db
		logger.Info(fmt.Sprintf("Preparing Postgres DB connection to %s...", cfg.PostgresDBHost))
		db, err = postgres.CreatePostgresDBConnection(cfg.PostgresDBHost, dbName, logger)
		if err != nil {
			return nil, fmt.Errorf("could not create postresql connection: %w", err)
		}
	}

	// Update historical transaction count from config if it's greater than 0
	if cfg.HistoricalTxCount > 0 {
		err = updateHistoricalCount(db, "historical_transaction_count", cfg.HistoricalTxCount)
		if err != nil {
			logger.Warn("Failed to update historical transaction count", "error", err)
		} else {
			logger.Info("Updated historical transaction count from config", "count", cfg.HistoricalTxCount)
		}
	}

	// Update historical contract count from config if it's greater than 0
	if cfg.HistoricalContractCount > 0 {
		err = updateHistoricalCount(db, "historical_contract_count", cfg.HistoricalContractCount)
		if err != nil {
			logger.Warn("Failed to update historical contract count", "error", err)
		} else {
			logger.Info("Updated historical contract count from config", "count", cfg.HistoricalContractCount)
		}
	}

	return hostdb.NewHostDB(db, logger)
}

// validateDBConf high-level checks that you have a valid configuration for DB creation
func validateDBConf(cfg *hostconfig.HostConfig) error {
	if cfg.UseInMemoryDB && cfg.PostgresDBHost != "" {
		return fmt.Errorf("invalid db config, useInMemoryDB=true so PostgresDB host not expected, but PostgresDBHost=%s", cfg.PostgresDBHost)
	}
	if cfg.SqliteDBPath != "" && cfg.UseInMemoryDB {
		return fmt.Errorf("useInMemoryDB=true so sqlite database will not be used and no path is needed, but sqliteDBPath=%s", cfg.SqliteDBPath)
	}
	return nil
}

// updateHistoricalCount updates the historical count for the specified table from config
func updateHistoricalCount(db *sqlx.DB, tableName string, count int) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET total = ? where id = 1", tableName)

	reboundQuery := db.Rebind(updateQuery)
	_, err := db.Exec(reboundQuery, count)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", tableName, err)
	}

	return nil
}
