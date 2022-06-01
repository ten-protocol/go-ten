package sql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/mattn/go-sqlite3" // this imports the sqlite driver to make the sql.Open() connection work
)

const (
	tempDirName = "temp-obscuro-persistence"
	createQry   = `create table if not exists kv (key binary(32) primary key, value blob); delete from kv;`
)

func CreateTemporarySQLiteDB(nodeID uint64) (ethdb.Database, error) {
	dbFile, err := getTempDBFile(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp sqlite DB file - %w", err)
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't create temp sqlite db - %w", err)
	}

	if _, err := db.Exec(createQry); err != nil {
		return nil, fmt.Errorf("failed to create sqlite db table - %w", err)
	}
	return CreateSQLEthDatabase(db)
}

func getTempDBFile(nodeID uint64) (string, error) {
	tempDir := filepath.Join("/tmp", tempDirName)
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create sqlite temp dir - %w", err)
	}
	// by using nodeIDs we ensure we overwrite old DBs when starting new tests
	tempFile := filepath.Join(tempDir, fmt.Sprintf("enclave-%v.db", nodeID))
	// delete old db if it exists
	_ = os.Remove(tempFile)
	return tempFile, nil
}
