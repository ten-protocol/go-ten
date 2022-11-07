package rawdb

import (
	"crypto/ecdsa"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

func ReadAttestationKey(db ethdb.KeyValueReader, address gethcommon.Address, logger gethlog.Logger) *ecdsa.PublicKey {
	key, err := db.Get(attestationPkKey(address))
	if err != nil {
		logger.Crit("Could not read key from db. ", log.ErrKey, err)
	}
	publicKey, err := crypto.DecompressPubkey(key)
	if err != nil {
		logger.Crit("Could not parse key from db.", log.ErrKey, err)
	}
	return publicKey
}

func WriteAttestationKey(db ethdb.KeyValueWriter, address gethcommon.Address, key *ecdsa.PublicKey, logger gethlog.Logger) {
	if err := db.Put(attestationPkKey(address), crypto.CompressPubkey(key)); err != nil {
		logger.Crit("Failed to store the attested key. ", log.ErrKey, err)
	}
}
