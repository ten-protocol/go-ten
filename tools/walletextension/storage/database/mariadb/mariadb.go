package mariadb

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/go-sql-driver/mysql" // Importing MariaDB driver
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database"
)

type MariaDB struct {
	db *sql.DB
}

// NewMariaDB creates a new MariaDB connection instance
func NewMariaDB(dbURL string) (*MariaDB, error) {
	db, err := sql.Open("mysql", dbURL+"?multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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

	return &MariaDB{db: db}, nil
}

func (m *MariaDB) AddUser(userID []byte, privateKey []byte) error {
	stmt, err := m.db.Prepare("REPLACE INTO users(user_id, private_key) VALUES (?, ?)")
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

func (m *MariaDB) DeleteUser(userID []byte) error {
	stmt, err := m.db.Prepare("DELETE FROM users WHERE user_id = ?")
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

func (m *MariaDB) GetUserPrivateKey(userID []byte) ([]byte, error) {
	var privateKey []byte
	err := m.db.QueryRow("SELECT private_key FROM users WHERE user_id = ?", userID).Scan(&privateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given userID
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return privateKey, nil
}

func (m *MariaDB) AddAccount(userID []byte, accountAddress []byte, signature []byte) error {
	stmt, err := m.db.Prepare("INSERT INTO accounts(user_id, account_address, signature) VALUES (?, ?, ?)")
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

func (m *MariaDB) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	rows, err := m.db.Query("SELECT account_address, signature FROM accounts WHERE user_id = ?", userID)
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

func (m *MariaDB) GetAllUsers() ([]common.UserDB, error) {
	rows, err := m.db.Query("SELECT user_id, private_key FROM users")
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
