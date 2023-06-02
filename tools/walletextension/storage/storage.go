package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	obscurocommon "github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpc"
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

func (s *Storage) SaveUserVK(userID string, vk *rpc.ViewingKey, message string) error {
	err := s.db.SaveUserVK(userID, vk, message)
	if err != nil {
		return fmt.Errorf("failed to save viewingkey to the storage, %w", err)
	}
	return nil
}

func (s *Storage) GetUserVKs(userID string) (map[common.Address]*rpc.ViewingKey, error) {
	userVKs, err := s.db.GetUserVKs(userID)
	if err != nil {
		return nil, err
	}
	return userVKs, nil
}

func (s *Storage) AddSignature(userID string, address []byte, signature string) error {
	err := s.db.AddSignature(userID, address, signature)
	if err != nil {
		fmt.Println("Error adding signature")
		return err
	}
	return nil
}

func (s *Storage) GetMessageAndSignature(userID string, address []byte) (string, string, error) {
	message, signature, err := s.db.GetMessageAndSignature(userID, address)
	if err != nil {
		fmt.Println("Error getting message and signature")
		return "", "", err
	}
	return message, signature, nil
}
