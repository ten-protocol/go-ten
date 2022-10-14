package db

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/config"
)

// CreateDBFromConfig creates an appropriate ethdb.Database instance based on your config
func CreateDBFromConfig(cfg config.EnclaveConfig, logger gethlog.Logger) (ethdb.Database, error) {
	if err := validateDBConf(cfg); err != nil {
		return nil, err
	}
	if cfg.UseInMemoryDB {
		logger.Info("UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		return getInMemDB()
	}

	if !cfg.WillAttest {
		// persistent but not secure in an enclave, we'll connect to a throwaway sqlite DB and test out persistence/sql implementations
		logger.Warn("Attestation is disabled, using a basic sqlite DB for persistence")
		// todo: for now we pass in an empty dbPath which will provide a throwaway temp file,
		// 		when we want to test persistence after node restart we should change this path to be config driven
		return sql.CreateTemporarySQLiteDB(cfg.SqliteDBPath, logger)
	}

	// persistent and with attestation means connecting to edgeless DB in a trusted enclave from a secure enclave
	logger.Info(fmt.Sprintf("Preparing Edgeless DB connection to %s...", cfg.EdgelessDBHost))
	return getEdgelessDB(cfg, logger)
}

// validateDBConf high-level checks that you have a valid configuration for DB creation
func validateDBConf(cfg config.EnclaveConfig) error {
	if cfg.UseInMemoryDB && cfg.EdgelessDBHost != "" {
		return fmt.Errorf("invalid db config, useInMemoryDB=true so EdgelessDB host not expected, but EdgelessDBHost=%s", cfg.EdgelessDBHost)
	}
	if !cfg.WillAttest && cfg.EdgelessDBHost != "" {
		return fmt.Errorf("invalid db config, willAttest=false so EdgelessDB host not supported, but EdgelessDBHost=%s", cfg.EdgelessDBHost)
	}
	if !cfg.UseInMemoryDB && cfg.WillAttest && cfg.EdgelessDBHost == "" {
		return fmt.Errorf("useInMemoryDB=false, willAttest=true so expected an EdgelessDB host but none was provided")
	}
	if cfg.SqliteDBPath != "" && cfg.UseInMemoryDB {
		return fmt.Errorf("useInMemoryDB=true so sqlite database will not be used and no path is needed, but sqliteDBPath=%s", cfg.SqliteDBPath)
	}
	if cfg.SqliteDBPath != "" && cfg.WillAttest {
		return fmt.Errorf("willAttest=true so sqlite database will not be used and no path is needed, but sqliteDBPath=%s", cfg.SqliteDBPath)
	}
	return nil
}

func getInMemDB() (ethdb.Database, error) {
	return rawdb.NewMemoryDatabase(), nil
}

func getEdgelessDB(cfg config.EnclaveConfig, logger gethlog.Logger) (ethdb.Database, error) {
	if cfg.EdgelessDBHost == "" {
		return nil, fmt.Errorf("failed to prepare EdgelessDB connection - EdgelessDBHost was not set on enclave config")
	}
	dbConfig := sql.EdgelessDBConfig{Host: cfg.EdgelessDBHost}
	return sql.EdgelessDBConnector(&dbConfig, logger)
}
