package core

import (
	"crypto/cipher"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

func DecryptTransactions(encryptedTxs nodecommon.EncryptedTransactions, rollupCipher cipher.AEAD) L2Txs {
	encodedTxs, err := rollupCipher.Open(nil, nil, encryptedTxs, nil)
	if err != nil {
		log.Panic("could not decrypt encrypted L2 transactions. Cause: %s", err)
	}

	txs := L2Txs{}
	if err := rlp.DecodeBytes(encodedTxs, txs); err != nil {
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
