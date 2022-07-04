package rawdb

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

func ReadAttestationKey(db ethdb.KeyValueReader, address gethcommon.Address) *ecdsa.PublicKey {
	key, err := db.Get(attestationPkKey(address))
	if err != nil {
		log.Panic("Could not read key from db. Cause: %s", err)
	}
	publicKey, err := crypto.DecompressPubkey(key)
	if err != nil {
		log.Panic("Could not parse key from db. Cause: %s", err)
	}
	return publicKey
}

func WriteAttestationKey(db ethdb.KeyValueWriter, address gethcommon.Address, key *ecdsa.PublicKey) {
	if err := db.Put(attestationPkKey(address), crypto.CompressPubkey(key)); err != nil {
		log.Panic("Failed to store the . Cause: %s", err)
	}
}
