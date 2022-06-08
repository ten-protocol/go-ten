package core

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// EnclavePrivateKeyHex is the ECDSA key used to securely communicate with the enclave.
// TODO - Replace this fixed key with a secret private key that unique per enclave.
const EnclavePrivateKeyHex = "2239e6a112b725625f355dcf7dd7a8d8a24a20f6448000a666b734f6957086a9"

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

// EncryptResponse - encrypts the response from the evm with the viewing key of the sender
func EncryptResponse(resp []byte) nodecommon.EncryptedResponse {
	// TODO
	return resp
}

// DecryptResponse - the reverse of EncryptResponse
func DecryptResponse(r nodecommon.EncryptedResponse) []byte {
	// TODO
	return r
}
