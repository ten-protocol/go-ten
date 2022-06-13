package core

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const (
	// TODO - This fixed nonce is insecure, and should be removed alongside the fixed rollup encryption key.
	rollupCipherNonce = "000000000000"
	// RollupEncryptionKeyHex is the AES key used to encrypt and decrypt the transaction blob in rollups.
	// TODO - Replace this fixed key with derived, rotating keys.
	RollupEncryptionKeyHex = "bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c"
)

// TransactionBlobCrypto handles the encryption and decryption of the transaction blobs stored inside a rollup.
type TransactionBlobCrypto interface {
	Encrypt(transactions L2Txs) nodecommon.EncryptedTransactions
	Decrypt(encryptedTxs nodecommon.EncryptedTransactions) L2Txs
}

type TransactionBlobCryptoImpl struct {
	transactionCipher cipher.AEAD
}

func NewTransactionBlobCryptoImpl() TransactionBlobCrypto {
	key := common.Hex2Bytes(RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("could not initialise AES cipher for enclave rollup key. Cause: %s", err))
	}
	transactionCipher, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("could not initialise wrapper for AES cipher for enclave rollup key. Cause: %s", err))
	}
	return TransactionBlobCryptoImpl{
		transactionCipher: transactionCipher,
	}
}

// TODO - Modify this logic so that transactions with different reveal periods are in different blobs, as per the whitepaper.
func (t TransactionBlobCryptoImpl) Encrypt(transactions L2Txs) nodecommon.EncryptedTransactions {
	encodedTxs, err := rlp.EncodeToBytes(transactions)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}

	return t.transactionCipher.Seal(nil, []byte(rollupCipherNonce), encodedTxs, nil)
}

func (t TransactionBlobCryptoImpl) Decrypt(encryptedTxs nodecommon.EncryptedTransactions) L2Txs {
	encodedTxs, err := t.transactionCipher.Open(nil, []byte(rollupCipherNonce), encryptedTxs, nil)
	if err != nil {
		log.Panic("could not decrypt encrypted L2 transactions. Cause: %s", err)
	}

	txs := L2Txs{}
	if err := rlp.DecodeBytes(encodedTxs, &txs); err != nil {
		log.Panic("could not decode encoded L2 transactions. Cause: %s", err)
	}

	return txs
}
