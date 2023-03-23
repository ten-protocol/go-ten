package sql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"

	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const (
	tempDirName = "obscuro-persistence"
	createQry   = `create table if not exists keyvalue (ky varbinary(64) primary key, val mediumblob);`
)

// CreateTemporarySQLiteDB if dbPath is empty will use a random throwaway temp file,
// otherwise dbPath is a filepath for the db file, allows for tests that care about persistence between restarts
func CreateTemporarySQLiteDB(dbPath string, logger gethlog.Logger) (*EnclaveDB, error) {
	if dbPath == "" {
		tempPath, err := CreateTempDBFile()
		if err != nil {
			return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
		}
		dbPath = tempPath
	}
	inMem := strings.Contains(dbPath, "mode=memory")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}

	// Sqlite in memory fails with table locks when there are multiple connections
	if inMem {
		db.SetMaxOpenConns(1)
	}

	err = initialiseDB(db)
	if err != nil {
		return nil, err
	}

	desc := "in memory"
	if !inMem {
		_, err := os.Stat(dbPath)
		if err == nil {
			desc = "existing"
		} else {
			desc = "new"
		}
	}
	logger.Info(fmt.Sprintf("Opened %s sqlite db file at %s", desc, dbPath))

	return CreateSQLEthDatabase(db, logger)
}

func initialiseDB(db *sql.DB) error {
	if _, err := db.Exec(createQry); err != nil {
		return fmt.Errorf("failed to create sqlite db table - %w", err)
	}
	return nil
}

func CreateTempDBFile() (string, error) {
	tempDir := filepath.Join("/tmp", tempDirName, randomStr(5))
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create sqlite temp dir - %w", err)
	}
	tempFile := filepath.Join(tempDir, "enclave.db")
	return tempFile, nil
}

// Generates a random string n characters long.
func randomStr(n int) string {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	suffix := make([]rune, n)
	for i := range suffix {
		suffix[i] = letters[randGen.Intn(len(letters))]
	}
	return string(suffix)
}
