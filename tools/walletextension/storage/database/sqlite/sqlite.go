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

	"github.com/ethereum/go-ethereum/crypto"

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

	return s.withTx(func(dbTx *sql.Tx) error {
		stmt, err := dbTx.Prepare("INSERT OR REPLACE INTO users(id, user_data) VALUES (?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(string(user.UserId), string(userJSON))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *SqliteDB) DeleteUser(userID []byte) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		stmt, err := dbTx.Prepare("DELETE FROM users WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(string(userID))
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		return nil
	})
}

func (s *SqliteDB) ActivateSessionKey(userID []byte, active bool) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		user, err := s.readUser(dbTx, userID)
		if err != nil {
			return err
		}
		user.ActiveSK = active
		return s.updateUser(dbTx, user)
	})
}

func (s *SqliteDB) AddSessionKey(userID []byte, key common.GWSessionKey) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		user, err := s.readUser(dbTx, userID)
		if err != nil {
			return err
		}
		user.SessionKey = &dbcommon.GWSessionKeyDB{
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

func (s *SqliteDB) RemoveSessionKey(userID []byte) error {
	return s.withTx(func(dbTx *sql.Tx) error {
		user, err := s.readUser(dbTx, userID)
		if err != nil {
			return err
		}
		user.SessionKey = nil
		return s.updateUser(dbTx, user)
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

func (s *SqliteDB) GetUser(userID []byte) (*common.GWUser, error) {
	var user dbcommon.GWUserDB
	var err error
	err = s.withTx(func(dbTx *sql.Tx) error {
		user, err = s.readUser(dbTx, userID)
		if err != nil {
			return err
		}
		return nil
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
			return dbcommon.GWUserDB{}, fmt.Errorf("failed to get user: %w", errutil.ErrNotFound)
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
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(updatedUserJSON), string(user.UserId))
	if err != nil {
		return fmt.Errorf("failed to update user with new account: %w", err)
	}

	return nil
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

func (s *SqliteDB) withTx(fn func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
