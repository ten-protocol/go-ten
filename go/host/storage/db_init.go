package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/host/storage/init/postgres"
)

const HOST = "HOST_"

// CreateDBFromConfig creates an appropriate ethdb.Database instance based on your config
func CreateDBFromConfig(cfg *config.HostConfig, logger gethlog.Logger) (hostdb.HostDB, error) {
	dbName := HOST + cfg.ID.String()
	if err := validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		sqliteDB, err := sqlite.CreateTemporarySQLiteHostDB(dbName, "mode=memory&cache=shared&_foreign_keys=on", "host_sqlite_init.sql")
		if err != nil {
			return nil, fmt.Errorf("could not create in memory sqlite DB: %w", err)
		}
		return hostdb.NewHostDB(sqliteDB, hostdb.SQLiteSQLStatements())
	}
	logger.Info(fmt.Sprintf("Preparing Postgres DB connection to %s...", cfg.PostgresDBHost))
	postgresDB, err := postgres.CreatePostgresDBConnection(cfg.PostgresDBHost, dbName, "host_postgres_init.sql")
	if err != nil {
		return nil, fmt.Errorf("could not create postresql connection: %w", err)
	}
	return hostdb.NewHostDB(postgresDB, hostdb.PostgresSQLStatements())
}

// validateDBConf high-level checks that you have a valid configuration for DB creation
func validateDBConf(cfg *config.HostConfig) error {
	if cfg.UseInMemoryDB && cfg.PostgresDBHost != "" {
		return fmt.Errorf("invalid db config, useInMemoryDB=true so MariaDB host not expected, but PostgresDBHost=%s", cfg.PostgresDBHost)
	}
	if cfg.SqliteDBPath != "" && cfg.UseInMemoryDB {
		return fmt.Errorf("useInMemoryDB=true so sqlite database will not be used and no path is needed, but sqliteDBPath=%s", cfg.SqliteDBPath)
	}
	return nil
}
