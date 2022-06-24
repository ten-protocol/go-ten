package core

import (
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// todo - joel - move to transaction injector

// EnclavePublicKeyHex is the public key of the enclave.
// TODO - Retrieve this key from the management contract instead.
const enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"

// EncryptTx encrypts a single transaction using the enclave's public key to send it privately to the enclave.
func EncryptTx(tx *nodecommon.L2Tx) (nodecommon.EncryptedTx, error) {
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("could not encode transaction bytes with RLP. Cause: %w", err)
	}

	// todo - joel - don't keep recomputing key, store somewhere.
	enclavePublicKey, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		return nil, fmt.Errorf("could not decompress enclave public key from hex. Cause: %w", err)
	}
	enclavePublicKeyEcies := ecies.ImportECDSAPublic(enclavePublicKey)

	encryptedTxBytes, err := ecies.Encrypt(rand.Reader, enclavePublicKeyEcies, txBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt request params with enclave public key. Cause: %w", err)
	}

	return encryptedTxBytes, nil
}
