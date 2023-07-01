package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/obscuronet/go-obscuro/tools/walletextension/common"

	obscurocommon "github.com/obscuronet/go-obscuro/go/common"
)

type Storage struct {
	db *SqliteDatabase
}

func New(dbPath string) (*Storage, error) {
	// If path is empty we create a random throwaway temp file, otherwise we use the path to the database
	if dbPath == "" {
		tempDir := filepath.Join("/tmp", "obscuro_gateway", obscurocommon.RandomStr(8))
		err := os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory: ", tempDir, err)
			return nil, err
		}
		dbPath = filepath.Join(tempDir, "gateway_databse.db")
	} else {
		dir := filepath.Dir(dbPath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			fmt.Println("Error creating directories:", err)
			return nil, err
		}
	}

	newDB, err := NewSqliteDatabase(dbPath)
	if err != nil {
		fmt.Println("Error creating database:", err)
		return nil, err
	}

	return &Storage{db: newDB}, nil
}

func (s *Storage) AddUser(userID []byte, privateKey []byte) error {
	err := s.db.AddUser(userID, privateKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteUser(userID []byte) error {
	err := s.db.DeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserPrivateKey(userID []byte) ([]byte, error) {
	privateKey, err := s.db.GetUserPrivateKey(userID)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func (s *Storage) AddAccount(userID []byte, accountAddress []byte, signature []byte) error {
	err := s.db.AddAccount(userID, accountAddress, signature)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *Storage) GetAllUsers() ([]common.UserDB, error) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
