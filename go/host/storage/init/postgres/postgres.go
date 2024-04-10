package postgres

import (
	"bufio"
	"database/sql"
	"embed"
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/lib/pq"
)

var (
	//go:embed *.sql
	sqlFiles embed.FS
)

const (
	defaultDatabase = "postgres"
	initFile        = "001_init.sql"
)

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
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	sqlFilePath := filepath.Join(baseDir, initFile)
	sqlFile, err := sqlFiles.ReadFile(sqlFilePath)

	if err != nil {
		return nil, fmt.Errorf("Could not read the initialisation sql file", log.ErrKey, err)
	}
	if err := initialiseDBFromSQLFile(db, string(sqlFile)); err != nil {
		return nil, fmt.Errorf("failed to initialize db from file %s: %w", sqlFile, err)
	}

	return db, nil
}

// initialiseDBFromSQLFile reads SQL commands from a file and executes them on the given DB connection.
func initialiseDBFromSQLFile(db *sql.DB, sqlFile string) error {
	file, err := os.Open(sqlFile)
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
