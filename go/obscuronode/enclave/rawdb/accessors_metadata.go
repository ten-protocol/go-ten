package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
)

func ReadSharedSecret(db ethdb.KeyValueReader) *core.SharedEnclaveSecret {
	var ss core.SharedEnclaveSecret

	enc, _ := db.Get(sharedSecret)
	if len(enc) == 0 {
		return nil
	}
	if err := rlp.DecodeBytes(enc, &ss); err != nil {
		return nil
	}

	return &ss
}

func WriteSharedSecret(db ethdb.KeyValueWriter, ss core.SharedEnclaveSecret) {
	enc, err := rlp.EncodeToBytes(ss)
	if err != nil {
		panic(err)
	}
	if err = db.Put(sharedSecret, enc); err != nil {
		panic(err)
	}
}

func ReadGenesisHash(db ethdb.KeyValueReader) *common.Hash {
	var hash common.Hash

	enc, _ := db.Get(genesisRollupHash)
	if len(enc) == 0 {
		return nil
	}
	if err := rlp.DecodeBytes(enc, &hash); err != nil {
		return nil
	}

	return &hash
}

func WriteGenesisHash(db ethdb.KeyValueWriter, hash common.Hash) {
	enc, err := rlp.EncodeToBytes(hash)
	if err != nil {
		panic(err)
	}
	if err = db.Put(genesisRollupHash, enc); err != nil {
		panic(err)
	}
}
