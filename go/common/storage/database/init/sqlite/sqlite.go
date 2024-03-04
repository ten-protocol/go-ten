package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common/storage/database/migration"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"

	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const (
	tempDirName = "ten-persistence"
)

//go:embed *.sql
var sqlFiles embed.FS

// CreateTemporarySQLiteEnclaveDB if dbPath is empty will use a random throwaway temp file,
// otherwise dbPath is a filepath for the sqldb file, allows for tests that care about persistence between restarts
func CreateTemporarySQLiteEnclaveDB(dbPath string, dbOptions string, logger gethlog.Logger, initFile string) (enclavedb.EnclaveDB, error) {
	initialsed := false

	if dbPath == "" {
		tempPath, err := CreateTempDBFile("enclave.db")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
		}
		dbPath = tempPath
	}

	inMem := strings.Contains(dbOptions, "mode=memory")
	description := "in memory"
	if !inMem {
		_, err := os.Stat(dbPath)
		if err == nil {
			description = "existing"
			initialsed = true
		} else {
			description = "new"
		}
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?%s", dbPath, dbOptions))
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}

	// Sqlite fails with table locks when there are multiple connections
	db.SetMaxOpenConns(1)

	if !initialsed {
		err = initialiseDB(db, initFile)
		if err != nil {
			return nil, err
		}
	}

	// perform db migration
	err = migration.CommonDBMigration(db, sqlFiles, logger.New(log.CmpKey, "DB_MIGRATION"))
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Opened %s sqlite db file at %s", description, dbPath))

	return enclavedb.NewEnclaveDB(db, logger)
}

// CreateTemporarySQLiteHostDB if dbPath is empty will use a random throwaway temp file,
// otherwise dbPath is a filepath for the sqldb file, allows for tests that care about persistence between restarts
func CreateTemporarySQLiteHostDB(dbPath string, dbOptions string, logger gethlog.Logger, initFile string) (*sql.DB, error) {
	if dbPath == "" {
		tempPath, err := CreateTempDBFile("host.db")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
		}
		dbPath = tempPath
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?%s", dbPath, dbOptions))
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}

	// Sqlite fails with table locks when there are multiple connections
	db.SetMaxOpenConns(1)

	err = initialiseDB(db, initFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise db - %w", err)
	}

	return db, nil
}

func initialiseDB(db *sql.DB, initFile string) error {
	sqlInitFile, err := sqlFiles.ReadFile(initFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlInitFile))
	if err != nil {
		return fmt.Errorf("failed to initialise sqlite %s - %w", sqlInitFile, err)
	}
	return nil
}

func CreateTempDBFile(dbname string) (string, error) {
	tempDir := filepath.Join("/tmp", tempDirName, common.RandomStr(5))
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create sqlite temp dir - %w", err)
	}
	tempFile := filepath.Join(tempDir, dbname)
	return tempFile, nil
}
