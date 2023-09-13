package storage

import (
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/storage/database"
)

type Storage interface {
	AddUser(userID []byte, privateKey []byte) error
	DeleteUser(userID []byte) error
	GetUserPrivateKey(userID []byte) ([]byte, error)
	AddAccount(userID []byte, accountAddress []byte, signature []byte) error
	GetAccounts(userID []byte) ([]common.AccountDB, error)
	GetAllUsers() ([]common.UserDB, error)
}

func New(dbPath string) (Storage, error) {
	return database.NewSqliteDatabase(dbPath)
}
