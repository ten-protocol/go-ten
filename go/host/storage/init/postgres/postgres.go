package postgres

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/lib/pq"

	"github.com/ten-protocol/go-ten/go/common/storage"

	_ "github.com/lib/pq"
)

const (
	defaultDatabase  = "postgres"
	maxDBConnections = 100
)

func CreatePostgresDBConnection(baseURL string, dbName string, logger gethlog.Logger) (*sql.DB, error) {
	driverName := registerPanicOnConnectionRefusedDriver(logger)
	if baseURL == "" {
		return nil, fmt.Errorf("failed to prepare PostgreSQL connection - DB URL was not set on host config")
	}
	dbURL := baseURL + defaultDatabase

	dbName = strings.ToLower(dbName)

	db, err := sql.Open(driverName, dbURL)
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
	db.SetMaxOpenConns(maxDBConnections)
	db.SetMaxIdleConns(maxDBConnections / 2)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

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

// registerPanicOnConnectionRefusedDriver registers the custom driver
func registerPanicOnConnectionRefusedDriver(logger gethlog.Logger) string {
	driverName := "pg-panic-on-unexpected-err"
	sql.Register(driverName,
		storage.NewPanicOnDBErrorDriver(
			&pq.Driver{},
			logger,
			func(err error) bool {
				return strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "shutting down")
			}),
	)
	return driverName
}
