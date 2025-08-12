package sqlite

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/storage/migration"

	"github.com/ten-protocol/go-ten/go/common"

	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const (
	tempDirName = "ten-persistence"
	initFile    = "001_init.sql"
)

//go:embed *.sql
var sqlFiles embed.FS

// CreateTemporarySQLiteHostDB if dbPath is empty will use a random throwaway temp file,
// otherwise dbPath is a filepath for the sqldb file, allows for tests that care about persistence between restarts
func CreateTemporarySQLiteHostDB(dbPath string, dbOptions string, logger gethlog.Logger) (*sqlx.DB, error) {
	if dbPath == "" {
		tempPath, err := CreateTempDBFile("host.db")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
		}
		dbPath = tempPath
	}

	db, err := sqlx.Open("sqlite3", fmt.Sprintf("file:%s?%s", dbPath, dbOptions))
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}

	// Sqlite fails with table locks when there are multiple connections
	db.SetMaxOpenConns(1)

	err = initialiseDB(db, initFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise db - %w", err)
	}

	err = migration.ApplyMigrations(db, sqlFiles, logger.New(log.CmpKey, "DB_MIGRATION"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initialiseDB(db *sqlx.DB, initFile string) error {
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
