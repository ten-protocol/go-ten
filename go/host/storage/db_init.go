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

type HostDB struct {
	DB    *sql.DB
	InMem bool
}

// CreateDBFromConfig creates an appropriate ethdb.Database instance based on your config
func CreateDBFromConfig(cfg *config.HostConfig, logger gethlog.Logger) (*HostDB, error) {
	if err := validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		println("CREATING IN MEM DB")
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		sqliteDb, err := sqlite.CreateTemporarySQLiteHostDB(HOST+cfg.ID.String(), "mode=memory&cache=shared&_foreign_keys=on", "host_sqlite_init.sql")
		if err != nil {
			return nil, fmt.Errorf("could not create in memory sqlite DB: %w", err)
		}
		return &HostDB{DB: sqliteDb, InMem: true}, nil
	}
	logger.Info(fmt.Sprintf("Preparing Maria DB connection to %s...", cfg.MariaDBHost))
	println("CREATING MARIA DB")
	mariaDb, err := mariadb.CreateMariaDBHostDB(cfg, HOST+cfg.ID.String(), "host_mariadb_init.sql")
	if err != nil {
		return nil, fmt.Errorf("could not create mariadb connection: %w", err)
	}
	return &HostDB{DB: mariaDb, InMem: false}, nil
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
