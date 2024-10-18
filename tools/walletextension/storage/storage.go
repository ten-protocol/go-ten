package storage

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/sqlite"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/encryption"
)

type Storage interface {
	AddUser(userID []byte, userIDHash []byte, privateKey []byte) error
	DeleteUser(userIDHash []byte) error
	GetUserPrivateKey(userIDHash []byte) ([]byte, error)
	AddAccount(userIDHash []byte, encryptedUserID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error
	GetAccounts(userIDHash []byte) ([]common.AccountDB, error)
	GetAllUsers() ([]common.UserDB, error)
	StoreTransaction(rawTx string, userIDHash []byte) error
}

type EncryptedStorage struct {
	Storage
	encryptor *encryption.Encryptor
}

func New(dbType string, dbConnectionURL, dbPath string, randomNumber []byte) (*EncryptedStorage, error) {
	var baseStorage Storage
	var err error

	switch dbType {
	// case "mariaDB":
	// 	baseStorage, err = mariadb.NewMariaDB(dbConnectionURL)
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
	userIDHash := es.encryptor.HashWithHMAC(userID)
	encryptedPrivateKey, err := es.encryptor.Encrypt(privateKey)
	if err != nil {
		return err
	}
	encryptedUserID, err := es.encryptor.Encrypt(userID)
	if err != nil {
		return err
	}
	return es.Storage.AddUser(encryptedUserID, userIDHash, encryptedPrivateKey)
}

func (es *EncryptedStorage) DeleteUser(userID []byte) error {
	userIDHash := es.encryptor.HashWithHMAC(userID)
	return es.Storage.DeleteUser(userIDHash)
}

func (es *EncryptedStorage) GetUserPrivateKey(userID []byte) ([]byte, error) {
	userIDHash := es.encryptor.HashWithHMAC(userID)
	encryptedPrivateKey, err := es.Storage.GetUserPrivateKey(userIDHash)
	if err != nil {
		return nil, err
	}
	return es.encryptor.Decrypt(encryptedPrivateKey)
}

func (es *EncryptedStorage) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error {
	userIDHash := es.encryptor.HashWithHMAC(userID)
	encryptedUserID, err := es.encryptor.Encrypt(userIDHash)
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

	return es.Storage.AddAccount(userIDHash, encryptedUserID, encryptedAccountAddress, encryptedSignature, signatureType)
}

func (es *EncryptedStorage) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	userIDHash := es.encryptor.HashWithHMAC(userID)
	accounts, err := es.Storage.GetAccounts(userIDHash)
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
		accounts[i].SignatureType = account.SignatureType
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

// StoreTransaction stores a transaction in the database
// TODO: Remove this method in the future as it is only a temporary solution for debugging purposes and that's why it's not encrypted
func (es *EncryptedStorage) StoreTransaction(rawTx string, userIDHash []byte) error {
	encryptedUserID := es.encryptor.HashWithHMAC(userIDHash)
	return es.Storage.StoreTransaction(rawTx, encryptedUserID)
}
