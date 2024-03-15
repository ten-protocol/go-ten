package storage

import (
	"database/sql"
	"fmt"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/host/storage/init/mariadb"
	"github.com/ten-protocol/go-ten/go/host/storage/init/sqlite"
)

const HOST = "HOST_"

// CreateDBFromConfig creates an appropriate ethdb.Database instance based on your config
func CreateDBFromConfig(cfg *config.HostConfig, logger gethlog.Logger) (*sql.DB, error) {
	if err := validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		return sqlite.CreateTemporarySQLiteHostDB(HOST+cfg.ID.String(), "mode=memory&cache=shared&_foreign_keys=on", "host_init.sql")
	}

	logger.Info(fmt.Sprintf("Preparing Maria DB connection to %s...", cfg.MariaDBHost))
	return mariadb.CreateMariaDBHostDB(cfg, "host_mariadb_init.sql")
}

// validateDBConf high-level checks that you have a valid configuration for DB creation
func validateDBConf(cfg *config.HostConfig) error {
	if cfg.UseInMemoryDB && cfg.MariaDBHost != "" {
		return fmt.Errorf("invalid db config, useInMemoryDB=true so MariaDB host not expected, but MariaDBHost=%s", cfg.MariaDBHost)
	}
	if cfg.SqliteDBPath != "" && cfg.UseInMemoryDB {
		return fmt.Errorf("useInMemoryDB=true so sqlite database will not be used and no path is needed, but sqliteDBPath=%s", cfg.SqliteDBPath)
	}
	return nil
}
