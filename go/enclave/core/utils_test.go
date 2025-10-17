package core

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
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

func TestGetAuthenticatedSenderWithZeroChainID(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	// pre-EIP155, no chain ID
	signer := types.HomesteadSigner{}
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	if _, err := GetAuthenticatedSender(testChainID, signedTx); err == nil {
		t.Errorf("expected error for transaction with zero chain ID, got nil")
	}
}

func TestGetAuthenticatedSenderWithValidChainID(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	expectedSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	sender, err := GetAuthenticatedSender(testChainID, signedTx)
	if err != nil {
		t.Fatalf("expected no error for valid transaction, got: %v", err)
	}
	if sender == nil {
		t.Fatal("expected sender address, got nil")
	}
	if *sender != expectedSender {
		t.Errorf("expected sender %s, got %s", expectedSender.Hex(), sender.Hex())
	}
}

func TestGetAuthenticatedSenderWithMismatchedChainID(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	wrongChainID := int64(testChainID + 1)
	signer := types.NewLondonSigner(big.NewInt(wrongChainID))
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	if _, err := GetAuthenticatedSender(testChainID, signedTx); err == nil {
		t.Errorf("expected error for mismatched chain ID, got nil")
	}
}

func TestGetExternalTxSignerWithValidChainID(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	expectedSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	sender, err := GetExternalTxSigner(signedTx)
	if err != nil {
		t.Fatalf("expected no error for valid transaction, got: %v", err)
	}
	if sender != expectedSender {
		t.Errorf("expected sender %s, got %s", expectedSender.Hex(), sender.Hex())
	}
}

func TestGetTxSignerWithZeroChainID(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	// pre-EIP155, no chain ID
	signer := types.HomesteadSigner{}
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	pricedTx := &common.L2PricedTransaction{
		Tx:             signedTx,
		PublishingCost: big.NewInt(0),
		FromSelf:       false,
		SystemDeployer: false,
	}

	if _, err := GetTxSigner(pricedTx); err == nil {
		t.Errorf("expected error for transaction with zero chain ID, got nil")
	}
}

func TestGetTxSignerWithValidChainID(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	pricedTx := &common.L2PricedTransaction{
		Tx:             signedTx,
		PublishingCost: big.NewInt(0),
		FromSelf:       false,
		SystemDeployer: false,
	}

	expectedSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	sender, err := GetTxSigner(pricedTx)
	if err != nil {
		t.Fatalf("expected no error for valid transaction, got: %v", err)
	}
	if sender != expectedSender {
		t.Errorf("expected sender %s, got %s", expectedSender.Hex(), sender.Hex())
	}
}

func TestGetTxSignerWithSystemDeployer(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	pricedTx := &common.L2PricedTransaction{
		Tx:             signedTx,
		PublishingCost: big.NewInt(0),
		FromSelf:       false,
		SystemDeployer: true,
	}

	// For SystemDeployer, the sender should be masked with the chain ID as address
	sender, err := GetTxSigner(pricedTx)
	if err != nil {
		t.Fatalf("expected no error for system deployer transaction, got: %v", err)
	}

	// Verify it's a masked sender (should not be the actual private key address)
	actualSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	if sender == actualSender {
		t.Errorf("expected masked sender for SystemDeployer, got actual sender %s", sender.Hex())
	}
}

func TestGetTxSignerWithFromSelf(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	txData := datagenerator.CreateL2TxData()
	signer := types.HomesteadSigner{}
	signedTx, _ := types.SignTx(types.NewTx(txData), signer, privateKey)

	pricedTx := &common.L2PricedTransaction{
		Tx:             signedTx,
		PublishingCost: big.NewInt(0),
		FromSelf:       true,
		SystemDeployer: false,
	}

	if signedTx.ChainId() != nil && signedTx.ChainId().Int64() != 0 {
		t.Fatalf("expected chain ID 0, got %v", signedTx.ChainId())
	}

	if _, err := GetTxSigner(pricedTx); err == nil {
		t.Errorf("expected error for transaction with zero chain ID, got nil")
	}
}
