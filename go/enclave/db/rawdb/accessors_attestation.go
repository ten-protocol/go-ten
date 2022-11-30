package rawdb

import (
	"crypto/ecdsa"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
)

func ReadAttestationKey(db ethdb.KeyValueReader, address gethcommon.Address) (*ecdsa.PublicKey, error) {
	key, err := db.Get(attestationPkKey(address))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := crypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func WriteAttestationKey(db ethdb.KeyValueWriter, address gethcommon.Address, key *ecdsa.PublicKey) error {
	if err := db.Put(attestationPkKey(address), crypto.CompressPubkey(key)); err != nil {
		return fmt.Errorf("could not write attestation key. Cause: %w", err)
	}
	return nil
}
