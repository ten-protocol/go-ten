package postgres

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/storage"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/lib/pq"
)

const (
	defaultDatabase = "postgres"
)

//go:embed *.sql
var sqlFiles embed.FS

func CreatePostgresDBConnection(baseURL string, dbName string) (*sql.DB, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("failed to prepare PostgreSQL connection - DB URL was not set on host config")
	}
	dbURL := baseURL + defaultDatabase

	dbName = strings.ToLower(dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL server: %v", err)
	}
	defer db.Close() // Close the connection when done

	rows, err := db.Query("SELECT 1 FROM pg_database WHERE datname = $1", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to query database existence: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database %s: %v", dbName, err)
		}
	}

	dbURL = fmt.Sprintf("%s%s", baseURL, dbName)

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database %s: %v", dbName, err)
	}

	// get the path to the migrations (they are always in the same directory as file containing connection function)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current directory")
	}
	migrationsDir := filepath.Dir(filename)

	if err = storage.ApplyMigrations(db, migrationsDir); err != nil {
		return nil, err
	}

	return db, nil
}
