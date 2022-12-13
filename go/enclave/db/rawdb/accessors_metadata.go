package rawdb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
)

func ReadSharedSecret(db ethdb.KeyValueReader) (*crypto.SharedEnclaveSecret, error) {
	var ss crypto.SharedEnclaveSecret

	enc, err := db.Get(sharedSecret)
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	if err := rlp.DecodeBytes(enc, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	return &ss, nil
}

func WriteSharedSecret(db ethdb.KeyValueWriter, ss crypto.SharedEnclaveSecret) error {
	enc, err := rlp.EncodeToBytes(ss)
	if err != nil {
		return fmt.Errorf("could not encode shared secret. Cause: %w", err)
	}
	if err = db.Put(sharedSecret, enc); err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	return nil
}
