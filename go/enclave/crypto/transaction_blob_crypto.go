package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/obscuronet/obscuro-playground/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"
)

const (
	// TODO - This fixed nonce is insecure, and should be removed alongside the fixed rollup encryption key.
	RollupCipherNonce = "000000000000"
	// RollupEncryptionKeyHex is the AES key used to encrypt and decrypt the transaction blob in rollups.
	// TODO - Replace this fixed key with derived, rotating keys.
	RollupEncryptionKeyHex = "bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c"
)

// TransactionBlobCrypto handles the encryption and decryption of the transaction blobs stored inside a rollup.
type TransactionBlobCrypto interface {
	Encrypt(transactions core.L2Txs) common.EncryptedTransactions
	Decrypt(encryptedTxs common.EncryptedTransactions) core.L2Txs
	// ToExtRollup - Transforms an internal rollup as seen by the enclave to an external rollup with an encrypted payload
	ToExtRollup(r *core.Rollup) common.ExtRollup
	ToEnclaveRollup(r *common.EncryptedRollup) *core.Rollup
}

type TransactionBlobCryptoImpl struct {
	transactionCipher cipher.AEAD
}

func NewTransactionBlobCryptoImpl() TransactionBlobCrypto {
	key := gethcommon.Hex2Bytes(RollupEncryptionKeyHex)
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
func (t TransactionBlobCryptoImpl) Encrypt(transactions core.L2Txs) common.EncryptedTransactions {
	encodedTxs, err := rlp.EncodeToBytes(transactions)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}

	return t.transactionCipher.Seal(nil, []byte(RollupCipherNonce), encodedTxs, nil)
}

func (t TransactionBlobCryptoImpl) Decrypt(encryptedTxs common.EncryptedTransactions) core.L2Txs {
	encodedTxs, err := t.transactionCipher.Open(nil, []byte(RollupCipherNonce), encryptedTxs, nil)
	if err != nil {
		log.Panic("could not decrypt encrypted L2 transactions. Cause: %s", err)
	}

	txs := core.L2Txs{}
	if err := rlp.DecodeBytes(encodedTxs, &txs); err != nil {
		log.Panic("could not decode encoded L2 transactions. Cause: %s", err)
	}

	return txs
}

func (t TransactionBlobCryptoImpl) ToExtRollup(r *core.Rollup) common.ExtRollup {
	return common.ExtRollup{
		Header: r.Header,
		// todo - joel - add txhashes here
		EncryptedTxBlob: t.Encrypt(r.Transactions),
	}
}

func (t TransactionBlobCryptoImpl) ToEnclaveRollup(r *common.EncryptedRollup) *core.Rollup {
	return &core.Rollup{
		Header:       r.Header,
		Transactions: t.Decrypt(r.Transactions),
	}
}
