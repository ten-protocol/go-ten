package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func ApplyMigrations(db *sql.DB, migrationsPath string) error {
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

func executeSQLFile(db *sql.DB, filePath string) error {
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
