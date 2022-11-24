package rawdb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
)

func ReadSharedSecret(db ethdb.KeyValueReader) (*crypto.SharedEnclaveSecret, error) {
	var ss crypto.SharedEnclaveSecret

	// TODO - Handle error.
	enc, _ := db.Get(sharedSecret)
	if len(enc) == 0 {
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

func ReadGenesisHash(db ethdb.KeyValueReader) (*common.Hash, bool) {
	var hash common.Hash

	enc, _ := db.Get(genesisRollupHash)
	if len(enc) == 0 {
		return nil, false
	}
	if err := rlp.DecodeBytes(enc, &hash); err != nil {
		return nil, false
	}

	return &hash, true
}

func WriteGenesisHash(db ethdb.KeyValueWriter, hash common.Hash, logger gethlog.Logger) {
	enc, err := rlp.EncodeToBytes(hash)
	if err != nil {
		logger.Crit("could not encode genesis hash. ", log.ErrKey, err)
	}
	if err = db.Put(genesisRollupHash, enc); err != nil {
		logger.Crit("could not put genesis hash in DB. ", log.ErrKey, err)
	}
}
