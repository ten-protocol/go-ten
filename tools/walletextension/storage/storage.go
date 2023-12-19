package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database"
)

type Storage interface {
	AddUser(userID []byte, privateKey []byte) error
	DeleteUser(userID []byte) error
	GetUserPrivateKey(userID []byte) ([]byte, error)
	AddAccount(userID []byte, accountAddress []byte, signature []byte) error
	GetAccounts(userID []byte) ([]common.AccountDB, error)
	GetAllUsers() ([]common.UserDB, error)
}

func New(dbType string, dbConnectionURL, dbPath string) (Storage, error) {
	switch dbType {
	case "mariaDB":
		return database.NewMariaDB(dbConnectionURL)
	case "sqlite":
		return database.NewSqliteDatabase(dbPath)
	}
	return nil, fmt.Errorf("unknown db %s", dbType)
}
