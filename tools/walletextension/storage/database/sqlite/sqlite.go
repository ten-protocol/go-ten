package sqlite

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ethereum/go-ethereum/crypto"

	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	obscurocommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
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

	// enable foreign keys in sqlite
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// create users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		user_id_hash BLOB PRIMARY KEY,
		user_id BLOB,
		private_key BLOB
	);`)
	if err != nil {
		return nil, err
	}

	// create accounts table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		user_id_hash BLOB,
		user_id BLOB,
		account_address BLOB,
		signature BLOB,
		signature_type int,
		FOREIGN KEY(user_id_hash) REFERENCES users(user_id_hash) ON DELETE CASCADE
	);`)
	if err != nil {
		return nil, err
	}

	// create transactions table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id BLOB,
    tx_hash TEXT,
    tx TEXT,
    tx_time TEXT DEFAULT (datetime('now'))
)	;`)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (s *Database) AddUser(userID []byte, userIDHash []byte, privateKey []byte) error {
	stmt, err := s.db.Prepare("INSERT OR REPLACE INTO users(user_id_hash, user_id, private_key) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userIDHash, userID, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) DeleteUser(userIDHash []byte) error {
	stmt, err := s.db.Prepare("DELETE FROM users WHERE user_id_hash = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userIDHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) GetUserPrivateKey(userIDhash []byte) ([]byte, error) {
	var privateKey []byte
	err := s.db.QueryRow("SELECT private_key FROM users WHERE user_id_hash = ?", userIDhash).Scan(&privateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given userIDHash
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return privateKey, nil
}

func (s *Database) AddAccount(userIDHash []byte, encryptedUserID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	stmt, err := s.db.Prepare("INSERT INTO accounts(user_id_hash, user_id, account_address, signature, signature_type) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userIDHash, encryptedUserID, accountAddress, signature, int(signatureType))
	if err != nil {
		return err
	}

	return nil
}

func (s *Database) GetAccounts(userIDHash []byte) ([]common.AccountDB, error) {
	rows, err := s.db.Query("SELECT account_address, signature, signature_type FROM accounts WHERE user_id_hash = ?", userIDHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []common.AccountDB
	for rows.Next() {
		var account common.AccountDB
		if err := rows.Scan(&account.AccountAddress, &account.Signature, &account.SignatureType); err != nil {
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

func (s *Database) StoreTransaction(rawTx string, userID []byte) error {
	stmt, err := s.db.Prepare("INSERT INTO transactions(user_id, tx_hash, tx) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	txHash := ""
	if len(rawTx) < 3 {
		fmt.Println("Invalid rawTx: ", rawTx)
	} else {
		// Decode the hex string to bytes, excluding the '0x' prefix
		rawTxBytes, err := hex.DecodeString(rawTx[2:])
		if err != nil {
			fmt.Println("Error decoding rawTx: ", err)
		} else {
			// Compute Keccak-256 hash
			txHash = crypto.Keccak256Hash(rawTxBytes).Hex()
		}
	}

	_, err = stmt.Exec(userID, txHash, rawTx)
	if err != nil {
		return err
	}

	return nil
}
