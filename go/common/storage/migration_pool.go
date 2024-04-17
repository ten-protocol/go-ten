package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ApplyMigrationsPool(pool *pgxpool.Pool, migrationsPath string) error {
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
		fmt.Println("Executing database migration file:", file)
		if err := executeSQLFilePool(pool, file); err != nil {
			return err
		}
	}

	return nil
}

func executeSQLFilePool(pool *pgxpool.Pool, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = pool.Exec(context.Background(), string(content))
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", filePath, err)
	}

	return nil
}
