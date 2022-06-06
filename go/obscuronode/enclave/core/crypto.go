package core

import (
	"crypto/cipher"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// TODO - This fixed nonce is insecure, and should be removed alongside the fixed rollup encryption key.
const rollupCipherNonce = "000000000000"

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

func EncryptTransactions(transactions L2Txs, rollupCipher cipher.AEAD) nodecommon.EncryptedTransactions {
	encodedTxs, err := rlp.EncodeToBytes(transactions)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}

	return rollupCipher.Seal(nil, []byte(rollupCipherNonce), encodedTxs, nil)
}

func DecryptTransactions(encryptedTxs nodecommon.EncryptedTransactions, rollupCipher cipher.AEAD) L2Txs {
	encodedTxs, err := rollupCipher.Open(nil, []byte(rollupCipherNonce), encryptedTxs, nil)
	if err != nil {
		log.Panic("could not decrypt encrypted L2 transactions. Cause: %s", err)
	}

	txs := L2Txs{}
	if err := rlp.DecodeBytes(encodedTxs, &txs); err != nil {
		log.Panic("could not decode encoded L2 transactions. Cause: %s", err)
	}

	return txs
}

func DecodeTx(tx nodecommon.EncodedTx) nodecommon.L2Tx {
	t := nodecommon.L2Tx{}
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		log.Panic("could not decrypt encrypted L2 transaction. Cause: %s", err)
	}

	return t
}

func EncodeTx(tx *nodecommon.L2Tx) nodecommon.EncodedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}
	return bytes
}
