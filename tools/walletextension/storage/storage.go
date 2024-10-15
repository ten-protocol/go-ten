package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/mariadb"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/sqlite"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/encryption"
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

type EncryptedStorage struct {
	Storage
	encryptor *encryption.Encryptor
}

func New(dbType string, dbConnectionURL, dbPath string, randomNumber []byte) (Storage, error) {
	var baseStorage Storage
	var err error

	switch dbType {
	case "mariaDB":
		baseStorage, err = mariadb.NewMariaDB(dbConnectionURL)
	case "sqlite":
		baseStorage, err = sqlite.NewSqliteDatabase(dbPath)
	default:
		return nil, fmt.Errorf("unknown db %s", dbType)
	}

	if err != nil {
		return nil, err
	}

	return NewEncryptedStorage(baseStorage, randomNumber)
}

func NewEncryptedStorage(storage Storage, key []byte) (*EncryptedStorage, error) {
	encryptor, err := encryption.NewEncryptor(key)
	if err != nil {
		return nil, err
	}
	return &EncryptedStorage{
		Storage:   storage,
		encryptor: encryptor,
	}, nil
}

// Implement the methods of EncryptedStorage
func (es *EncryptedStorage) AddUser(userID []byte, privateKey []byte) error {
	encryptedPrivateKey, err := es.encryptor.Encrypt(privateKey)
	if err != nil {
		return err
	}
	return es.Storage.AddUser(userID, encryptedPrivateKey)
}

func (es *EncryptedStorage) GetUserPrivateKey(userID []byte) ([]byte, error) {
	encryptedPrivateKey, err := es.Storage.GetUserPrivateKey(userID)
	if err != nil {
		return nil, err
	}
	return es.encryptor.Decrypt(encryptedPrivateKey)
}
func (es *EncryptedStorage) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	encryptedUserID, err := es.encryptor.Encrypt(userID)
	if err != nil {
		return err
	}
	encryptedAccountAddress, err := es.encryptor.Encrypt(accountAddress)
	if err != nil {
		return err
	}
	encryptedSignature, err := es.encryptor.Encrypt(signature)
	if err != nil {
		return err
	}

	return es.Storage.AddAccount(encryptedUserID, encryptedAccountAddress, encryptedSignature, signatureType)
}

func (es *EncryptedStorage) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	accounts, err := es.Storage.GetAccounts(userID)
	if err != nil {
		return nil, err
	}

	for i, account := range accounts {
		decryptedAccountAddress, err := es.encryptor.Decrypt(account.AccountAddress)
		if err != nil {
			return nil, err
		}
		decryptedSignature, err := es.encryptor.Decrypt(account.Signature)
		if err != nil {
			return nil, err
		}

		accounts[i].AccountAddress = decryptedAccountAddress
		accounts[i].Signature = decryptedSignature
	}

	return accounts, nil
}

func (es *EncryptedStorage) GetAllUsers() ([]common.UserDB, error) {
	users, err := es.Storage.GetAllUsers()
	if err != nil {
		return nil, err
	}

	decryptedUsers := make([]common.UserDB, len(users))
	for i, user := range users {
		decryptedUserID, err := es.encryptor.Decrypt(user.UserID)
		if err != nil {
			return nil, err
		}
		decryptedPrivateKey, err := es.encryptor.Decrypt(user.PrivateKey)
		if err != nil {
			return nil, err
		}
		decryptedUsers[i] = common.UserDB{
			UserID:     decryptedUserID,
			PrivateKey: decryptedPrivateKey,
		}
	}
	return decryptedUsers, nil
}

func (es *EncryptedStorage) StoreTransaction(rawTx string, userID []byte) error {
	encryptedUserID, err := es.encryptor.Encrypt(userID)
	if err != nil {
		return err
	}

	encryptedRawTx, err := es.encryptor.Encrypt([]byte(rawTx))
	if err != nil {
		return err
	}

	return es.Storage.StoreTransaction(string(encryptedRawTx), encryptedUserID)
}
