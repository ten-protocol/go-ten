package sqlite

/*
	SQLite database implementation of the Storage interface

	SQLite is used for local deployments and testing without the need for a cloud database.
	To make sure to see similar behaviour as in production using CosmosDB we use SQLite database in a similar way as comosDB (as key-value database).
*/
import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	dbcommon "github.com/ten-protocol/go-ten/tools/walletextension/storage/database/common"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"

	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	obscurocommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

type SqliteDB struct {
	db *sql.DB
}

const sqliteCfg = "_foreign_keys=on&_journal_mode=wal&_txlock=immediate&_synchronous=normal"

func NewSqliteDatabase(dbPath string) (*SqliteDB, error) {
	// load the db file
	dbFilePath, err := createOrLoad(dbPath)
	if err != nil {
		return nil, err
	}

	// open the db
	path := fmt.Sprintf("file:%s?%s", dbFilePath, sqliteCfg)
	db, err := sql.Open("sqlite3", path)
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

	return &SqliteDB{db: db}, nil
}

func (s *SqliteDB) AddUser(userID []byte, privateKey []byte) error {
	user := dbcommon.GWUserDB{
		UserId:     userID,
		PrivateKey: privateKey,
		Accounts:   []dbcommon.GWAccountDB{},
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

	_, err = stmt.Exec(string(user.UserId), string(userJSON))
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteDB) DeleteUser(userID []byte) error {
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

func (s *SqliteDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	var userDataJSON string
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	var user dbcommon.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	newAccount := dbcommon.GWAccountDB{
		AccountAddress: accountAddress,
		Signature:      signature,
		SignatureType:  int(signatureType),
	}

	user.Accounts = append(user.Accounts, newAccount)

	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling updated user: %w", err)
	}

	stmt, err := tx.Prepare("UPDATE users SET user_data = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(updatedUserJSON), string(userID))
	if err != nil {
		return fmt.Errorf("failed to update user with new account: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteDB) GetUser(userID []byte) (*common.GWUser, error) {
	var userDataJSON string
	err := s.db.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get user: %w", errutil.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user dbcommon.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user.ToGWUser(), nil
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
