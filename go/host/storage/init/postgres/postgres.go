package postgres

import (
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/storage/migration"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/lib/pq"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/storage"
)

const (
	defaultDatabase  = "postgres"
	maxDBConnections = 75 // azure has 100 max connections
	initFile         = "001_init.sql"
)

//go:embed *.sql
var sqlFiles embed.FS

func CreatePostgresDBConnection(baseURL string, dbName string, logger gethlog.Logger) (*sqlx.DB, error) {
	driverName := registerPanicOnConnectionRefusedDriver(logger)
	if baseURL == "" {
		return nil, fmt.Errorf("failed to prepare PostgreSQL connection - DB URL was not set on host config")
	}
	dbURL := baseURL + defaultDatabase

	dbName = strings.ToLower(dbName)

	// Open connection to the default postgres database to check if our target DB exists
	defaultDB, err := sqlx.Open(driverName, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL server: %v", err)
	}
	defer defaultDB.Close() // Close the default postgres connection when done

	rows, err := defaultDB.Query("SELECT 1 FROM pg_database WHERE datname = $1", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to query database existence: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database %s: %v", dbName, err)
		}
	}

	// close the default postgres connection explicitly before opening the target DB
	defaultDB.Close()

	dbURL = fmt.Sprintf("%s%s", baseURL, dbName)

	// open conneciton to target DB
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database %s: %v", dbName, err)
	}
	db.SetMaxOpenConns(maxDBConnections)
	db.SetMaxIdleConns(maxDBConnections / 3)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	// initialise the database with the initial SQL file
	err = migration.InitialiseDB(db, sqlFiles, initFile)
	if err != nil {
		return nil, err
	}

	// apply any additional migrations
	err = migration.ApplyMigrations(db, sqlFiles, logger.New(log.CmpKey, "DB_MIGRATION"))
	if err != nil {
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
			},
			func(conn driver.Conn) {
			}),
	)

	// tell sqlx this driver uses PostgreSQL syntax ($1, $2, $3)
	sqlx.BindDriver(driverName, sqlx.DOLLAR)

	logger.Info("Registered custom PostgreSQL driver with panic handling", "driver_name", driverName)
	return driverName
}
