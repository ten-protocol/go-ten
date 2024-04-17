package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ten-protocol/go-ten/go/common/storage"
)

const (
	defaultDatabase       = "postgres"
	maxConnections        = 100
	maxConnectionLifetime = time.Minute * 30
)

func CreatePostgresDBConnection(baseURL string, dbName string) (*pgxpool.Pool, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("failed to prepare PostgreSQL connection - DB URL was not set on host config")
	}
	dbURL := baseURL + defaultDatabase

	dbName = strings.ToLower(dbName)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %v", err)
	}

	poolConfig.MaxConns = maxConnections
	poolConfig.MaxConnLifetime = maxConnectionLifetime

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL server: %v", err)
	}

	// Ensure the connection pool is closed when done
	defer pool.Close()

	_, err = pool.Exec(context.Background(), "CREATE DATABASE IF NOT EXISTS "+dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to create database %s: %v", dbName, err)
	}

	dbURL = fmt.Sprintf("%s%s", baseURL, dbName)

	poolConfig, err = pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %v", err)
	}

	pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database %s: %v", dbName, err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current directory")
	}
	migrationsDir := filepath.Dir(filename)

	if err = storage.ApplyMigrationsPool(pool, migrationsDir); err != nil {
		return nil, err
	}

	return pool, nil
}
