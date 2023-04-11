package rawdb

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
)

var _enclaveKeyDBKey = []byte("ek")

func StoreEnclaveKey(db ethdb.KeyValueWriter, enclaveKey *ecdsa.PrivateKey, _ gethlog.Logger) error {
	if enclaveKey == nil {
		return errors.New("enclaveKey cannot be nil")
	}
	keyBytes := crypto.FromECDSA(enclaveKey)
	return db.Put(_enclaveKeyDBKey, keyBytes)
}

func GetEnclaveKey(db ethdb.KeyValueReader, _ gethlog.Logger) (*ecdsa.PrivateKey, error) {
	keyBytes, err := db.Get(_enclaveKeyDBKey)
	if err != nil {
		return nil, err
	}
	enclaveKey, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to construct ECDSA private key from enclave key bytes - %w", err)
	}
	return enclaveKey, nil
}
