package mariadb

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"os"
	"strings"
)

func CreateMariaDBHostDB(cfg *config.HostConfig, initFile string) (*sql.DB, error) {
	mariaDbHost := cfg.MariaDBHost
	if mariaDbHost == "" {
		return nil, fmt.Errorf("failed to prepare MariaDB connection - MariaDBHost was not set on host config")
	}
	db, err := sql.Open("mysql", mariaDbHost+"?multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if initFile != "" {
		if err := initialiseDBFromSQLFile(db, initFile); err != nil {
			return nil, fmt.Errorf("failed to initialise db from file %s: %w", initFile, err)
		}
	}

	//_, filename, _, ok := runtime.Caller(0)
	//if !ok {
	//	return nil, fmt.Errorf("failed to get current directory")
	//}
	//migrationsDir := filepath.Dir(filename)
	//if err = database.ApplyMigrations(db, migrationsDir); err != nil {
	//	return nil, err
	//}

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
