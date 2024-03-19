package mariadb

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	maxDBPoolSize = 100
)

func CreateMariaDBHostDB(dbURL string, dbName string, initFile string) (*sql.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("failed to prepare MariaDB connection - MariaDBHost was not set on host config")
	}
	db, err := sql.Open("mysql", dbURL+"?multiStatements=true")
	if err != nil {
		log.Fatalf("Failed to connect to MariaDB server: %v", err)
	}
	db.SetMaxOpenConns(maxDBPoolSize)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	_, err = db.Exec(fmt.Sprintf("USE `%s`", dbName))
	if err != nil {
		log.Fatalf("Failed to select database %s: %v", dbName, err)
	}

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	sqlFile := filepath.Join(baseDir, "host_mariadb_init.sql")

	if initFile != "" {
		if err := initialiseDBFromSQLFile(db, sqlFile); err != nil {
			println("Error initialisting DB from file")
			return nil, fmt.Errorf("failed to initialise db from file %s: %w", initFile, err)
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
