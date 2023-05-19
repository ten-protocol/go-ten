package storage

import (
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	"github.com/obscuronet/go-obscuro/go/rpc"
)

type SqliteDatabase struct {
	db *sql.DB
}

func NewSqliteDatabase(dbName string) (*SqliteDatabase, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS viewingkeys (
		user_id binary(32),
		account_address binary(20),
		private_key binary(32),
		signed_key binary(65)
	);`)

	if err != nil {
		return nil, err
	}

	return &SqliteDatabase{db: db}, nil
}

func (s *SqliteDatabase) SaveUserVK(userID string, vk *rpc.ViewingKey) error {
	stmt, err := s.db.Prepare("INSERT INTO viewingkeys (user_id, account_address, private_key, signed_key) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error creating sql statement ", err)
		return err
	}
	defer stmt.Close()

	viewingPrivateKeyBytes := crypto.FromECDSA(vk.PrivateKey.ExportECDSA())
	_, err = stmt.Exec(userID, vk.Account.Bytes(), viewingPrivateKeyBytes, vk.SignedKey)
	if err != nil {
		fmt.Println("Error inserting failed")
		return err
	}

	return nil
}

func (s *SqliteDatabase) GetUserVKs(userID string) (map[common.Address]*rpc.ViewingKey, error) {
	viewingKeys := make(map[common.Address]*rpc.ViewingKey)

	rows, err := s.db.Query("SELECT user_id, account_address, private_key, signed_key FROM viewingkeys WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("Error in getting items from db", err)
		return nil, err
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		fmt.Println("Error in getting rows from db", err)
		return nil, err
	}

	var tmpUserID []byte
	var tmpAccountAddress []byte
	var tmpPrivateKey []byte
	var tmpSignedKey []byte
	for rows.Next() {
		err := rows.Scan(&tmpUserID, &tmpAccountAddress, &tmpPrivateKey, &tmpSignedKey)
		if err != nil {
			fmt.Println("Error in looping over results")
			return nil, err
		}

		account := common.BytesToAddress(tmpAccountAddress)

		viewingKeyPrivate, err := crypto.ToECDSA(tmpPrivateKey)
		if err != nil {
			fmt.Println("Error ToECDSA conversion", err)
			return nil, err
		}

		viewingKeys[account] = &rpc.ViewingKey{
			Account:    &account,
			PrivateKey: ecies.ImportECDSA(viewingKeyPrivate),
			PublicKey:  crypto.CompressPubkey(&viewingKeyPrivate.PublicKey),
			SignedKey:  tmpSignedKey,
		}
	}

	return viewingKeys, nil
}
