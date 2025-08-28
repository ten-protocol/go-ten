package sqlite

/*
	SQLite database implementation of the Storage interface.

	This implementation mimics the CosmosDB approach where we store the entire user record (including accounts and session keys)
	in a single JSON object within the 'users' table. There are no separate tables for accounts or session keys.

	Each user record:
	{
		"user_data": <entire GWUserDB JSON>
	}

	This simplifies the schema and keeps it similar to the CosmosDB container-based storage.
*/

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()

	dbcommon "github.com/ten-protocol/go-ten/tools/walletextension/storage/database/common"

	obscurocommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type SqliteDB struct {
	db *sql.DB
}

const sqliteCfg = "_foreign_keys=on&_journal_mode=wal&_txlock=immediate&_synchronous=normal"

func NewSqliteDatabase(dbPath string) (*SqliteDB, error) {
	// load or create the db file
	dbFilePath, err := createOrLoad(dbPath)
	if err != nil {
		return nil, err
	}

	// open the db
	path := fmt.Sprintf("file:%s?%s", dbFilePath, sqliteCfg)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Enable foreign keys in SQLite (harmless, even though we don't use them now)
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// Create the users table if it doesn't exist. We store entire user as JSON.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		user_data TEXT
	);`)
	if err != nil {
		return nil, err
	}

	// If there was an old 'accounts' table from a previous implementation, drop it.
	// This ensures no leftover foreign key constraints cause issues.
	_, _ = db.Exec("DROP TABLE IF EXISTS accounts;")

	return &SqliteDB{db: db}, nil
}

func (s *SqliteDB) AddUser(userID []byte, privateKey []byte) error {
	user := dbcommon.GWUserDB{
		UserId:      userID,
		PrivateKey:  privateKey,
		Accounts:    []dbcommon.GWAccountDB{},
		SessionKeys: make(map[common.Address]*dbcommon.GWSessionKeyDB),
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	return s.withTx(func(dbTx *sql.Tx) error {
		stmt, err := dbTx.Prepare("INSERT OR REPLACE INTO users(id, user_data) VALUES (?, ?)")
		if err != nil {
			return fmt.Errorf("failed to prepare insert statement: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(string(user.UserId), string(userJSON))
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
		return nil
	})
}

func (s *SqliteDB) DeleteUser(userID []byte) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		stmt, err := dbTx.Prepare("DELETE FROM users WHERE id = ?")
		if err != nil {
			return fmt.Errorf("failed to prepare delete statement: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(string(userID))
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return nil
	})
}

func (s *SqliteDB) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		user, err := s.readUser(dbTx, userID)
		if err != nil {
			return err
		}

		newAccount := dbcommon.GWAccountDB{
			AccountAddress: accountAddress,
			Signature:      signature,
			SignatureType:  int(signatureType),
		}

		user.Accounts = append(user.Accounts, newAccount)
		return s.updateUser(dbTx, user)
	})
}

func (s *SqliteDB) AddSessionKey(userID []byte, key wecommon.GWSessionKey) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		user, err := s.readUser(dbTx, userID)
		if err != nil {
			return err
		}

		// Check session key limit
		if len(user.SessionKeys) >= wecommon.MaxSessionKeysPerUser {
			return fmt.Errorf("maximum number of session keys (%d) reached", wecommon.MaxSessionKeysPerUser)
		}

		// Initialize SessionKeys map if nil
		if user.SessionKeys == nil {
			user.SessionKeys = make(map[common.Address]*dbcommon.GWSessionKeyDB)
		}

		address := *key.Account.Address
		user.SessionKeys[address] = &dbcommon.GWSessionKeyDB{
			PrivateKey: crypto.FromECDSA(key.PrivateKey.ExportECDSA()),
			Account: dbcommon.GWAccountDB{
				AccountAddress: key.Account.Address.Bytes(),
				Signature:      key.Account.Signature,
				SignatureType:  int(key.Account.SignatureType),
			},
		}
		return s.updateUser(dbTx, user)
	})
}

func (s *SqliteDB) RemoveSessionKey(userID []byte, sessionKeyAddr *common.Address) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		user, err := s.readUser(dbTx, userID)
		if err != nil {
			return err
		}

		if user.SessionKeys == nil {
			return fmt.Errorf("no session keys found for user")
		}

		if _, exists := user.SessionKeys[*sessionKeyAddr]; !exists {
			return fmt.Errorf("session key not found: %s", sessionKeyAddr.Hex())
		}

		delete(user.SessionKeys, *sessionKeyAddr)
		return s.updateUser(dbTx, user)
	})
}

func (s *SqliteDB) GetUser(userID []byte) (*wecommon.GWUser, error) {
	var user dbcommon.GWUserDB
	var err error
	err = s.withTx(func(dbTx *sql.Tx) error {
		user, err = s.readUser(dbTx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return user.ToGWUser()
}

func (s *SqliteDB) readUser(dbTx *sql.Tx, userID []byte) (dbcommon.GWUserDB, error) {
	var userDataJSON string
	err := dbTx.QueryRow("SELECT user_data FROM users WHERE id = ?", string(userID)).Scan(&userDataJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbcommon.GWUserDB{}, dbcommon.ErrUserNotFound
		}
		return dbcommon.GWUserDB{}, fmt.Errorf("failed to get user: %w", err)
	}

	var user dbcommon.GWUserDB
	err = json.Unmarshal([]byte(userDataJSON), &user)
	if err != nil {
		return dbcommon.GWUserDB{}, fmt.Errorf("failed to unmarshal user data: %w", err)
	}
	return user, nil
}

func (s *SqliteDB) updateUser(dbTx *sql.Tx, user dbcommon.GWUserDB) error {
	updatedUserJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling updated user: %w", err)
	}

	stmt, err := dbTx.Prepare("UPDATE users SET user_data = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(updatedUserJSON), string(user.UserId))
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// GetEncryptionKey returns nil for SQLite as it doesn't use encryption directly in this implementation.
func (s *SqliteDB) GetEncryptionKey() []byte {
	return nil
}

func (s *SqliteDB) withTx(fn func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func createOrLoad(dbPath string) (string, error) {
	// If path is empty we create a random temporary file, otherwise we use the provided path
	if dbPath == "" {
		tempDir := filepath.Join("/tmp", "obscuro_gateway", obscurocommon.RandomStr(8))
		err := os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error creating directory %s: %w", tempDir, err)
		}
		dbPath = filepath.Join(tempDir, "gateway_database.db")
	} else {
		dir := filepath.Dir(dbPath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return "", fmt.Errorf("error creating directories: %w", err)
		}
	}

	return dbPath, nil
}
