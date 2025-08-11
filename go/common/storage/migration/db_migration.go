package migration

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ten-protocol/go-ten/go/common/storage"

	"github.com/jmoiron/sqlx"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

const currentMigrationVersionKey = "CURRENT_MIGRATION_VERSION"

func DBMigration(db *sqlx.DB, sqlFiles embed.FS, logger gethlog.Logger) error {
	migrationFiles, err := readMigrationFiles(sqlFiles)
	if err != nil {
		return err
	}

	maxMigration := int64(len(migrationFiles))

	var maxDB int64
	config, err := storage.FetchConfig(context.Background(), db, currentMigrationVersionKey)
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

// InitialiseDB initializes the database with the initial SQL file
func InitialiseDB(db *sqlx.DB, sqlFiles embed.FS, initFile string) error {
	sqlInitFile, err := sqlFiles.ReadFile(initFile)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to initialise postgres: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(string(sqlInitFile))
	if err != nil {
		return fmt.Errorf("failed to initialise postgres %s: %w", sqlInitFile, err)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func executeMigration(db *sqlx.DB, content string, migrationOrder int64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Split statements by semicolon and execute each one
	statements := strings.Split(content, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue // Skip empty statements
		}

		_, err = tx.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s - %w", stmt, err)
		}
	}

	err = storage.InsertOrUpdateConfig(context.Background(), tx, currentMigrationVersionKey, big.NewInt(migrationOrder).Bytes())
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

func ApplyMigrationFromPath(db *sqlx.DB, migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	var sqlFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsPath, file.Name()))
		}
	}

	sort.Strings(sqlFiles) // Sort files lexicographically to apply migrations in order

	for _, file := range sqlFiles {
		fmt.Println("Executing db migration file: ", file)
		if err = executeSQLFile(db, file); err != nil {
			return err
		}
	}

	return nil
}

func executeSQLFile(db *sqlx.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", filePath, err)
	}

	return nil
}
