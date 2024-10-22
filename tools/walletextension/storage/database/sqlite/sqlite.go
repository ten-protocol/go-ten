package sqlite

/*
	SQLite database implementation of the Storage interface

	SQLite is used for local deployments and testing without the need for a cloud database.
	To make sure to see similar behaviour as in production using CosmosDB we use SQLite database in a similar way as comosDB (as key-value database).
*/
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	obscurocommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
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

	// enable foreign keys in sqlite
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// Modify the users table to store the entire GWUserDB as JSON
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		user_data TEXT
	);`)
	if err != nil {
		return nil, err
	}

	// Remove the accounts table as it will be stored within the user_data JSON

	return &Database{db: db}, nil
}

func (s *Database) AddUser(userID []byte, privateKey []byte) error {
	user := common.GWUserDB{
		ID:         string(userID),
		UserId:     userID,
		PrivateKey: privateKey,
		Accounts:   []common.GWAccountDB{},
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	stmt, err := s.db.Prepare("INSERT OR REPLACE INTO users(id, user_data) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, string(userJSON))
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) DeleteUser(userID []byte) error {
	stmt, err := s.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(userID))
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *Database) GetUserPrivateKey(userID []byte) ([]byte, error) {
	var userDataJSON string
	err := s.db.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	var user common.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user.PrivateKey, nil
}

func (s *Database) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	var userDataJSON string
	err := s.db.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	var user common.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	newAccount := common.GWAccountDB{
		AccountAddress: accountAddress,
		Signature:      signature,
		SignatureType:  int(signatureType),
	}

	user.Accounts = append(user.Accounts, newAccount)

	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling updated user: %w", err)
	}

	stmt, err := s.db.Prepare("UPDATE users SET user_data = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(updatedUserJSON), string(userID))
	if err != nil {
		return fmt.Errorf("failed to update user with new account: %w", err)
	}

	return nil
}

func (s *Database) GetAccounts(userID []byte) ([]common.GWAccountDB, error) {
	var userDataJSON string
	err := s.db.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user common.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user.Accounts, nil
}

func (s *Database) GetUser(userID []byte) (common.GWUserDB, error) {
	var userDataJSON string
	err := s.db.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.GWUserDB{}, fmt.Errorf("failed to get user: %w", errutil.ErrNotFound)
		}
		return common.GWUserDB{}, fmt.Errorf("failed to get user: %w", err)
	}

	var user common.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return common.GWUserDB{}, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user, nil
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
