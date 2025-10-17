package core

import (
	"math/big"
	"testing"

	"github.com/ten-protocol/go-ten/integration/datagenerator"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const testChainID = 1234

func TestValidSignatureVerifies(t *testing.T) {
	tx := datagenerator.CreateL2Tx()
	privateKey, _ := crypto.GenerateKey()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	signedTx, _ := types.SignTx(tx, signer, privateKey)

	if err := VerifySignature(testChainID, signedTx); err != nil {
		t.Errorf("validly-signed transaction did not pass verification: %v", err)
	}
}

func TestUnsignedTxDoesNotVerify(t *testing.T) {
	tx := datagenerator.CreateL2Tx()

	if err := VerifySignature(testChainID, tx); err == nil {
		t.Errorf("transaction was not signed but verified anyway: %v", err)
	}
}

func TestModifiedTxDoesNotVerify(t *testing.T) {
	txData := datagenerator.CreateL2TxData()
	tx := types.NewTx(txData)
	privateKey, _ := crypto.GenerateKey()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	_, _ = types.SignTx(tx, signer, privateKey)

	// We create a new transaction around the transaction data, breaking the signature.
	modifiedTx := types.NewTx(txData)

	if err := VerifySignature(testChainID, modifiedTx); err == nil {
		t.Errorf("transaction was modified after signature but verified anyway: %v", err)
	}
}

func TestIncorrectSignerDoesNotVerify(t *testing.T) {
	tx := datagenerator.CreateL2Tx()
	privateKey, _ := crypto.GenerateKey()
	incorrectChainID := int64(testChainID + 1)
	signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
	signedTx, _ := types.SignTx(tx, signer, privateKey)

	if err := VerifySignature(testChainID, signedTx); err == nil {
		t.Errorf("transaction used incorrect signer but verified anyway: %v", err)
	}
}

func TestZeroChainIDDoesNotPanic(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	// pre-EIP155, no chain ID
	signer := types.HomesteadSigner{}
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	if signedTx.ChainId() != nil && signedTx.ChainId().Int64() != 0 {
		t.Fatalf("expected chain ID 0, got %v", signedTx.ChainId())
	}
	
	if _, err := GetExternalTxSigner(signedTx); err == nil {
		t.Errorf("expected error for transaction with zero chain ID, got nil")
	}
}
