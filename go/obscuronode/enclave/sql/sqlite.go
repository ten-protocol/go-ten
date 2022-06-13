package sql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const (
	tempDirName = "obscuro-persistence"
	createQry   = `create table if not exists kv (key binary(32) primary key, value blob); delete from kv;`
)

// CreateTemporarySQLiteDB if dbPath is empty will use a random throwaway temp file,
// 	otherwise dbPath is a filepath for the db file, allows for tests that care about persistence between restarts
func CreateTemporarySQLiteDB(dbPath string) (ethdb.Database, error) {
	if dbPath == "" {
		tempPath, err := getTempDBFile()
		if err != nil {
			return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
		}
		dbPath = tempPath
	}
	// determine if a db file already exists, we don't want to overwrite it
	_, err := os.Stat(dbPath)
	existingDB := err == nil

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open sqlite db - %w", err)
	}
	desc := "existing"
	if !existingDB {
		// db wasn't there already so we should set it up (create kv store table)
		if _, err := db.Exec(createQry); err != nil {
			return nil, fmt.Errorf("failed to create sqlite db table - %w", err)
		}
		desc = "new"
	}
	log.Info("Opened %s sqlite db file at %s", desc, dbPath)
	return CreateSQLEthDatabase(db)
}

func getTempDBFile() (string, error) {
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
