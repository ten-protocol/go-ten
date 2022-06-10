package core

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

// EncryptTx encrypts a single transaction using the enclave's public key to send it privately to the enclave.
// TODO - Perform real encryption here, and not just RLP encoding.
func EncryptTx(tx *nodecommon.L2Tx) nodecommon.EncryptedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}
	return bytes
}

// DecryptTx reverses the encryption performed by EncryptTx.
// TODO - Perform real decryption here, and not just RLP decoding.
func DecryptTx(tx nodecommon.EncryptedTx) nodecommon.L2Tx {
	t := nodecommon.L2Tx{}
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		log.Panic("could not decrypt encrypted L2 transaction. Cause: %s", err)
	}

	return t
}
