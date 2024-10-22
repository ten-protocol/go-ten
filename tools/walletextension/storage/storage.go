package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/sqlite"
)

type Storage interface {
	AddUser(userID []byte, privateKey []byte) error
	DeleteUser(userID []byte) error
	AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error
	GetUser(userID []byte) (common.GWUserDB, error)
}

func New(dbType string, dbConnectionURL, dbPath string) (Storage, error) {
	// TODO @ziga: Generate random key in a different part of the code!
	randomKey, err := common.GenerateRandomKey(32)
	if err != nil {
		return nil, err
	}

	switch dbType {
	case "sqlite":
		return sqlite.NewSqliteDatabase(dbPath)
	case "cosmosDB":
		return cosmosdb.NewCosmosDB(dbConnectionURL, randomKey)
	}

	return nil, fmt.Errorf("unknown db %s", dbType)
}
