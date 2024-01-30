package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	obscurocommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database"

	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	common "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type Database struct {
	db *sql.DB
}

func NewSqliteDatabase(dbPath string) (*Database, error) {
	// load the db file
	dbFilePath, err := createOrLoad(dbPath)
	if err != nil {
		return nil, err
	}

	// open the db
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return nil, err
	}

	// get the path to the migrations (they are always in the same directory as file containing connection function)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current directory")
	}
	migrationsDir := filepath.Dir(filename)

	// apply migrations
	if err = database.ApplyMigrations(db, migrationsDir); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (s *Database) AddUser(userID []byte, privateKey []byte) error {
	stmt, err := s.db.Prepare("INSERT OR REPLACE INTO users(user_id, private_key) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) DeleteUser(userID []byte) error {
	stmt, err := s.db.Prepare("DELETE FROM users WHERE user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) GetUserPrivateKey(userID []byte) ([]byte, error) {
	var privateKey []byte
	err := s.db.QueryRow("SELECT private_key FROM users WHERE user_id = ?", userID).Scan(&privateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given userID
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return privateKey, nil
}

func (s *Database) AddAccount(userID []byte, accountAddress []byte, signature []byte) error {
	stmt, err := s.db.Prepare("INSERT INTO accounts(user_id, account_address, signature) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, accountAddress, signature)
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	rows, err := s.db.Query("SELECT account_address, signature FROM accounts WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []common.AccountDB
	for rows.Next() {
		var account common.AccountDB
		if err := rows.Scan(&account.AccountAddress, &account.Signature); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Database) GetAllUsers() ([]common.UserDB, error) {
	rows, err := s.db.Query("SELECT user_id, private_key FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []common.UserDB
	for rows.Next() {
		var user common.UserDB
		err = rows.Scan(&user.UserID, &user.PrivateKey)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func createOrLoad(dbPath string) (string, error) {
	// If path is empty we create a random throwaway temp file, otherwise we use the path to the database
	if dbPath == "" {
		tempDir := filepath.Join("/tmp", "obscuro_gateway", obscurocommon.RandomStr(8))
		err := os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory: ", tempDir, err)
			return "", err
		}
		dbPath = filepath.Join(tempDir, "gateway_databse.db")
	} else {
		dir := filepath.Dir(dbPath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			fmt.Println("Error creating directories:", err)
			return "", err
		}
	}

	return dbPath, nil
}
