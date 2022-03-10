package enclave

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestValidSignatureVerifies(t *testing.T) {
	tx := createL2Tx()
	privateKey, _ := crypto.GenerateKey()
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	signedTx, _ := types.SignTx(tx, signer, privateKey)

	if err := verifySignature(signedTx); err != nil {
		t.Errorf("validly-signed transaction did not pass verification: %s", err)
	}
}

func TestUnsignedTxDoesNotVerify(t *testing.T) {
	tx := createL2Tx()

	if err := verifySignature(tx); err == nil {
		t.Errorf("transaction was not signed but verified anyway: %s", err)
	}
}

func TestModifiedTxDoesNotVerify(t *testing.T) {
	txData := createL2TxData()
	tx := types.NewTx(txData)
	privateKey, _ := crypto.GenerateKey()
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	_, _ = types.SignTx(tx, signer, privateKey)

	// We modify the transaction's nonce after signing, breaking the signature.
	txData.Nonce = 0
	modifiedTx := types.NewTx(txData)

	if err := verifySignature(modifiedTx); err == nil {
		t.Errorf("transaction was modified after signature but verified anyway: %s", err)
	}
}

func TestIncorrectSignerDoesNotVerify(t *testing.T) {
	tx := createL2Tx()
	privateKey, _ := crypto.GenerateKey()
	incorrectChainID := int64(ChainID + 1)
	signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
	signedTx, _ := types.SignTx(tx, signer, privateKey)

	if err := verifySignature(signedTx); err == nil {
		t.Errorf("transaction used incorrect signer but verified anyway: %s", err)
	}
}
