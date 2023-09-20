package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"sort"
	"strconv"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/storage/enclavedb"
)

const currentMigrationVersionKey = "CURRENT_MIGRATION_VERSION"

func DBMigration(db *sql.DB, sqlFiles embed.FS, logger gethlog.Logger) error {
	migrationFiles, err := readMigrationFiles(sqlFiles)
	if err != nil {
		return err
	}

	maxMigration := int64(len(migrationFiles))

	var maxDB int64
	config, err := enclavedb.FetchConfig(db, currentMigrationVersionKey)
	if err != nil {
		// first time there is no entry, so 001 was executed already ( triggered at launch/manifest time )
		if errors.Is(err, errutil.ErrNotFound) {
			maxDB = 1
		} else {
			return err
		}
	} else {
		maxDB = ByteArrayToInt(config)
	}

	// write to the database
	for i := maxDB; i < maxMigration; i++ {
		logger.Info("Executing db migration", "file", migrationFiles[i].Name())
		content, err := sqlFiles.ReadFile(migrationFiles[i].Name())
		if err != nil {
			return err
		}
		err = executeMigration(db, string(content), i)
		if err != nil {
			return fmt.Errorf("unable to execute migration for %s - %w", migrationFiles[i].Name(), err)
		}
		logger.Info("Successfully executed", "file", migrationFiles[i].Name(), "index", i)
	}

	return nil
}

func executeMigration(db *sql.DB, content string, migrationOrder int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(content)
	if err != nil {
		return err
	}

	_, err = enclavedb.WriteConfigToTx(tx, currentMigrationVersionKey, big.NewInt(migrationOrder).Bytes())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func readMigrationFiles(sqlFiles embed.FS) ([]fs.DirEntry, error) {
	migrationFiles, err := sqlFiles.ReadDir(".")
	if err != nil {
		return nil, err
	}

	// sort the migrationFiles based on the prefix (before "_")
	sort.Slice(migrationFiles, func(i, j int) bool {
		// Extract the number prefix and compare
		return getPrefix(migrationFiles[i]) < getPrefix(migrationFiles[j])
	})

	// sanity check the consecutive rule
	for i, file := range migrationFiles {
		prefix := getPrefix(file)
		if i+1 != prefix {
			panic("Invalid migration file. Missing index")
		}
	}
	return migrationFiles, err
}

func getPrefix(migrationFile fs.DirEntry) int {
	prefix := strings.Split(migrationFile.Name(), "_")[0]
	number, err := strconv.Atoi(prefix)
	if err != nil {
		panic("Invalid db migration file")
	}
	return number
}

func ByteArrayToInt(arr []byte) int64 {
	b := big.NewInt(0)
	b.SetBytes(arr)
	return b.Int64()
}
