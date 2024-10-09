package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/mariadb"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/sqlite"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type Storage interface {
	AddUser(userID []byte, privateKey []byte) error
	DeleteUser(userID []byte) error
	GetUserPrivateKey(userID []byte) ([]byte, error)
	AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error
	GetAccounts(userID []byte) ([]common.AccountDB, error)
	GetAllUsers() ([]common.UserDB, error)
	StoreTransaction(rawTx string, userID []byte) error
}

func New(dbType string, dbConnectionURL, dbPath string) (Storage, error) {
	switch dbType {
	case "mariaDB":
		return mariadb.NewMariaDB(dbConnectionURL)
	case "sqlite":
		return sqlite.NewSqliteDatabase(dbPath)
	case "cosmosDB":
		return cosmosdb.NewCosmosDB(dbConnectionURL)
	}
	return nil, fmt.Errorf("unknown db %s", dbType)
}
