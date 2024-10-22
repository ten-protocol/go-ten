package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

type Storage interface {
	AddUser(userID []byte, privateKey []byte) error
	DeleteUser(userID []byte) error
	GetUserPrivateKey(userID []byte) ([]byte, error)
	AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error
	GetAccounts(userID []byte) ([]common.AccountDB, error)
	GetUser(userID []byte) (common.UserDB, error)
}

func New(dbType string, dbConnectionURL, dbPath string) (Storage, error) {
	switch dbType {
	// TODO: @ziga - fix sqlite for new Storage interface to allow local testing and running it locally without cosmosdb
	// case "sqlite":
	// 	return sqlite.NewSqliteDatabase(dbPath)
	case "cosmosDB":
		return cosmosdb.NewCosmosDB(dbConnectionURL)
	}

	return nil, fmt.Errorf("unknown db %s", dbType)
}
