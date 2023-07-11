package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	common "github.com/obscuronet/go-obscuro/tools/walletextension/common"
)

type SqliteDatabase struct {
	db *sql.DB
}

func NewSqliteDatabase(dbName string) (*SqliteDatabase, error) {
	db, err := sql.Open("sqlite3", dbName)
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
		user_id binary(32) PRIMARY KEY,
		private_key binary(32)
	);`)

	if err != nil {
		return nil, err
	}

	// create accounts table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		user_id binary(32),
		account_address binary(20),
		signature binary(65),
    	FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`)

	if err != nil {
		return nil, err
	}

	return &SqliteDatabase{db: db}, nil
}

func (s *SqliteDatabase) AddUser(userID []byte, privateKey []byte) error {
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

func (s *SqliteDatabase) DeleteUser(userID []byte) error {
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

func (s *SqliteDatabase) GetUserPrivateKey(userID []byte) ([]byte, error) {
	var privateKey []byte
	err := s.db.QueryRow("SELECT private_key FROM users WHERE user_id = ?", userID).Scan(&privateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given userID
			return nil, nil
		}
		return nil, err
	}

	return privateKey, nil
}

func (s *SqliteDatabase) AddAccount(userID []byte, accountAddress []byte, signature []byte) error {
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

func (s *SqliteDatabase) GetAccounts(userID []byte) ([]common.AccountDB, error) {
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

func (s *SqliteDatabase) GetAllUsers() ([]common.UserDB, error) {
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
