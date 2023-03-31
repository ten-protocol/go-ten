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
	tempDirName       = "obscuro-persistence"
	createKVTable     = `create table if not exists keyvalue (ky varbinary(64) primary key, val mediumblob);`
	createEventsTable = "create table if not exists events (" +
		" topic0 binary(32), " +
		" topic1 binary(32), " +
		" topic2 binary(32), " +
		" topic3 binary(32), " +
		" topic4 binary(32), " +
		" datablob mediumblob, " +
		" blockHash binary(32), " +
		" blockNumber int, " +
		" txHash binary(32), " +
		" txIdx int, " +
		" logIdx int, " +
		" address binary(32), " +
		" lifecycleEvent boolean, " +
		" relAddress1 binary(20), " +
		" relAddress2 binary(20), " +
		" relAddress3 binary(20), " +
		" relAddress4 binary(20) " +
		");" +
		"create index IX_AD on events(address);" +
		"create index IX_RAD1 on events(relAddress1);" +
		"create index IX_RAD2 on events(relAddress2);" +
		"create index IX_RAD3 on events(relAddress3);" +
		"create index IX_RAD4 on events(relAddress4);" +
		"create index IX_BLH on events(blockHash);" +
		"create index IX_BLN on events(blockNumber);" +
		"create index IX_TXH on events(txHash);" +
		"create index IX_T0 on events(topic0);" +
		"create index IX_T1 on events(topic1);" +
		"create index IX_T2 on events(topic2);" +
		"create index IX_T3 on events(topic3);" +
		"create index IX_T4 on events(topic4);"
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
	description := "in memory"
	if !inMem {
		_, err := os.Stat(dbPath)
		if err == nil {
			description = "existing"
		} else {
			description = "new"
		}
	}

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

	logger.Info(fmt.Sprintf("Opened %s sqlite db file at %s", description, dbPath))

	return CreateSQLEthDatabase(db, logger)
}

func initialiseDB(db *sql.DB) error {
	if _, err := db.Exec(createKVTable); err != nil {
		return fmt.Errorf("failed to create sqlite kv db table - %w", err)
	}
	if _, err := db.Exec(createEventsTable); err != nil {
		return fmt.Errorf("failed to create sqlite events db table - %w", err)
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
