package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ten-protocol/go-ten/go/common/log"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/storage/init/migration"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"

	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const (
	tempDirName = "ten-persistence"
	initFile    = "001_init.sql"
)

//go:embed *.sql
var sqlFiles embed.FS

// CreateTemporarySQLiteDB if dbPath is empty will use a random throwaway temp file,
// otherwise dbPath is a filepath for the sqldb file, allows for tests that care about persistence between restarts
// We create 2 sqlite instances. One R/W with a single connection, and a R/O with multiple connections
func CreateTemporarySQLiteDB(dbPath string, dbOptions string, config enclaveconfig.EnclaveConfig, logger gethlog.Logger) (enclavedb.EnclaveDB, error) {
	initialsed := false

	if dbPath == "" {
		tempPath, err := CreateTempDBFile()
		if err != nil {
			return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
		}
		dbPath = tempPath
	}

	var description string

	_, err := os.Stat(dbPath)
	if err == nil {
		description = "existing"
		initialsed = true
	} else {
		myfile, e := os.Create(dbPath)
		if e != nil {
			logger.Crit("could not create temp sqlite DB file - %w", e)
		}
		myfile.Close()

		description = "new"
	}

	path := fmt.Sprintf("file:%s?mode=rw&%s", dbPath, dbOptions)
	logger.Info("Connect to sqlite", "path", path)
	rwdb, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}

	// Sqlite fails with table locks when there are multiple connections
	rwdb.SetMaxOpenConns(1)

	if !initialsed {
		err = initialiseDB(rwdb)
		if err != nil {
			return nil, err
		}
	}

	// perform db migration
	err = migration.DBMigration(rwdb, sqlFiles, logger.New(log.CmpKey, "DB_MIGRATION"))
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Opened %s sqlite db file at %s", description, dbPath))

	roPath := fmt.Sprintf("file:%s?mode=ro&%s", dbPath, dbOptions)
	logger.Info("Connect to sqlite", "ro_path", roPath)
	rodb, err := sql.Open("sqlite3", roPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}
	rodb.SetMaxOpenConns(10)

	return enclavedb.NewEnclaveDB(rodb, rwdb, config, logger)
}

func initialiseDB(db *sql.DB) error {
	sqlInitFile, err := sqlFiles.ReadFile(initFile)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to initialise sqlite  %w", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(string(sqlInitFile))
	if err != nil {
		return fmt.Errorf("failed to initialise sqlite %s - %w", sqlInitFile, err)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func CreateTempDBFile() (string, error) {
	tempDir := filepath.Join("/tmp", tempDirName, common.RandomStr(5))
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create sqlite temp dir - %w", err)
	}
	tempFile := filepath.Join(tempDir, "enclave.db")
	return tempFile, nil
}
