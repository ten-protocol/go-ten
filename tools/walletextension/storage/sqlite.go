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

	// todo (@ziga) - should we use binary format also for signature and text
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS viewingkeys (
		user_id binary(32),
		account_address binary(20),
		private_key binary(32),
		signed_key binary(65),
    	message	TEXT,
    	message_signature TEXT,
    	last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		return nil, err
	}

	return &SqliteDatabase{db: db}, nil
}

func (s *SqliteDatabase) SaveUserVK(userID string, vk *rpc.ViewingKey, message string) error {
	stmt, err := s.db.Prepare("INSERT INTO viewingkeys (user_id, account_address, private_key, signed_key, message) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error creating sql statement ", err)
		return err
	}
	defer stmt.Close()

	viewingPrivateKeyBytes := crypto.FromECDSA(vk.PrivateKey.ExportECDSA())
	_, err = stmt.Exec(userID, vk.Account.Bytes(), viewingPrivateKeyBytes, vk.SignedKey, message)
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

func (s *SqliteDatabase) GetMessageAndSignature(userID string, address []byte) (string, string, error) {
	var message string
	var messageSignature sql.NullString

	// Execute the SQL statement with the provided parameters
	row := s.db.QueryRow("SELECT message, message_signature FROM viewingkeys WHERE user_id = ? AND account_address = ?", userID, address)
	err := row.Scan(&message, &messageSignature)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no rows were found
			fmt.Println("No rows found for the given userID and address.")
			return "", "", nil
		}
		fmt.Println("Failed to retrieve message and message_signature:", err)
		return "", "", err
	}

	if messageSignature.Valid {
		return message, messageSignature.String, nil
	}
	return message, "", nil
}

func (s *SqliteDatabase) AddSignature(userID string, address []byte, signature string) error {
	// Execute the SQL statement with the provided parameters
	result, err := s.db.Exec("UPDATE viewingkeys SET message_signature = ? WHERE user_id = ? AND account_address = ?", signature, userID, address)
	if err != nil {
		fmt.Println("Failed to update message_signature:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("No rows were affected by the update.")
		// Handle the case where zero rows were affected
		// Return an appropriate error
	}

	return nil
}
