package storage

import (
	"fmt"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/host/storage/init/postgres"
)

const HOST = "HOST_"
const sqliteHostCfg = "_foreign_keys=on&_journal_mode=wal&_txlock=immediate&_synchronous=normal&mode=memory&cache=shared"

// CreateDBFromConfig creates an appropriate ethdb.Database instance based on your config
func CreateDBFromConfig(cfg *hostconfig.HostConfig, logger gethlog.Logger) (hostdb.HostDB, error) {
	dbName := HOST + cfg.ID.String()
	if err := validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		sqliteDB, err := sqlite.CreateTemporarySQLiteHostDB(dbName, sqliteHostCfg)
		if err != nil {
			return nil, fmt.Errorf("could not create in memory sqlite DB: %w", err)
		}
		return hostdb.NewHostDB(sqliteDB, hostdb.SQLiteSQLStatements())
	}
	logger.Info(fmt.Sprintf("Preparing Postgres DB connection to %s...", cfg.PostgresDBHost))
	postgresDB, err := postgres.CreatePostgresDBConnection(cfg.PostgresDBHost, dbName)
	if err != nil {
		return nil, fmt.Errorf("could not create postresql connection: %w", err)
	}
	return hostdb.NewHostDB(postgresDB, hostdb.PostgresSQLStatements())
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
