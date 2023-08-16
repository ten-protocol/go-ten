package sqlite

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/obscuronet/go-obscuro/go/enclave/storage/enclavedb"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"

	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const tempDirName = "obscuro-persistence"

//go:embed 001_init.sql
var sqlInitFile string

// CreateTemporarySQLiteDB if dbPath is empty will use a random throwaway temp file,
// otherwise dbPath is a filepath for the sqldb file, allows for tests that care about persistence between restarts
func CreateTemporarySQLiteDB(dbPath string, dbOptions string, logger gethlog.Logger) (enclavedb.EnclaveDB, error) {
	initialsed := false

	if dbPath == "" {
		tempPath, err := CreateTempDBFile()
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

	db.SetMaxOpenConns(1)

	if !initialsed {
		err = initialiseDB(db)
		if err != nil {
			return nil, err
		}
	}

	logger.Info(fmt.Sprintf("Opened %s sqlite db file at %s", description, dbPath))

	return enclavedb.NewEnclaveDB(db, logger)
}

func initialiseDB(db *sql.DB) error {
	_, err := db.Exec(sqlInitFile)
	if err != nil {
		return fmt.Errorf("failed to initialise sqlite %s - %w", sqlInitFile, err)
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
