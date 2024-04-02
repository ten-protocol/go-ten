package postgres

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/lib/pq"
)

const (
	maxDBPoolSize = 100
)

func CreatePostgresDBConnection(dbURL string, dbName string, initFile string) (*sql.DB, error) {
	//if dbURL == "" {
	//	return nil, fmt.Errorf("failed to prepare PostgreSQL connection - DB URL was not set on host config")
	//}
	dbURL = "postgres://WillHester:1866@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL server: %v", err)
	}
	db.SetMaxOpenConns(maxDBPoolSize)
	// Check if the database exists
	rows, err := db.Query("SELECT 1 FROM pg_database WHERE datname = $1", dbName)
	if err != nil {
		log.Fatalf("Failed to query database existence: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		// Database doesn't exist, create it
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("Failed to create database %s: %v", dbName, err)
		}
	}

	//_, err = db.Exec(fmt.Sprintf("SET DATABASE %s", dbName))
	//if err != nil {
	//	log.Fatalf("Failed to select database %s: %v", dbName, err)
	//}

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	sqlFile := filepath.Join(baseDir, "host_postgres_init.sql")

	if initFile != "" {
		if err := initialiseDBFromSQLFile(db, sqlFile); err != nil {
			println("Error initializing DB from file")
			return nil, fmt.Errorf("failed to initialize db from file %s: %w", initFile, err)
		}
	}

	return db, nil
}

// initialiseDBFromSQLFile reads SQL commands from a file and executes them on the given DB connection.
func initialiseDBFromSQLFile(db *sql.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sqlStatement string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "--") || strings.TrimSpace(line) == "" {
			continue
		}
		sqlStatement += line
		if strings.HasSuffix(line, ";") {
			_, err := db.Exec(sqlStatement)
			if err != nil {
				return fmt.Errorf("failed to execute SQL statement: %s, error: %w", sqlStatement, err)
			}
			sqlStatement = "" // Reset statement
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading SQL file: %w", err)
	}

	return nil
}
