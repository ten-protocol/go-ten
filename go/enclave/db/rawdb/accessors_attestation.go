package rawdb

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
)

func ReadAttestationKey(db ethdb.KeyValueReader, address gethcommon.Address) (*ecdsa.PublicKey, error) {
	key, err := db.Get(attestationPkKey(address))
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := crypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func WriteAttestationKey(db ethdb.KeyValueWriter, address gethcommon.Address, key *ecdsa.PublicKey, logger gethlog.Logger) error {
	if err := db.Put(attestationPkKey(address), crypto.CompressPubkey(key)); err != nil {
		return fmt.Errorf("could not write attestation key. Cause: %w", err)
	}
	return nil
}
