package enclave

import (
	"math/big"
	"testing"

	"github.com/obscuronet/obscuro-playground/integration/datagenerator"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const testChainID = 1234

func TestValidSignatureVerifies(t *testing.T) {
	tx := datagenerator.CreateL2Tx()
	privateKey, _ := crypto.GenerateKey()
	signer := types.NewLondonSigner(big.NewInt(testChainID))
	signedTx, _ := types.SignTx(tx, signer, privateKey)

	if err := verifySignature(testChainID, signedTx); err != nil {
		t.Errorf("validly-signed transaction did not pass verification: %v", err)
	}
}

func TestUnsignedTxDoesNotVerify(t *testing.T) {
	tx := datagenerator.CreateL2Tx()

	if err := verifySignature(testChainID, tx); err == nil {
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

	if err := verifySignature(testChainID, modifiedTx); err == nil {
		t.Errorf("transaction was modified after signature but verified anyway: %v", err)
	}
}

func TestIncorrectSignerDoesNotVerify(t *testing.T) {
	tx := datagenerator.CreateL2Tx()
	privateKey, _ := crypto.GenerateKey()
	incorrectChainID := int64(testChainID + 1)
	signer := types.NewLondonSigner(big.NewInt(incorrectChainID))
	signedTx, _ := types.SignTx(tx, signer, privateKey)

	if err := verifySignature(testChainID, signedTx); err == nil {
		t.Errorf("transaction used incorrect signer but verified anyway: %v", err)
	}
}

func TestInvalidBlocksAreRejected(t *testing.T) {
	// There are no tests of acceptance of valid chains of blocks. This is because the logic to generate a valid block
	// is non-trivial.
	genesisJSON, err := core.DefaultGenesisBlock().MarshalJSON()
	if err != nil {
		t.Errorf("could not parse genesis JSON: %v", err)
	}
	enclave := enclaveImpl{l1Blockchain: NewL1Blockchain(genesisJSON)}

	invalidHeaders := []types.Header{
		{ParentHash: common.HexToHash("0x0")},                                                            // Unknown ancestor.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(999)},            // Wrong block number.
		{ParentHash: core.DefaultGenesisBlock().ToBlock(nil).Hash(), Number: big.NewInt(1), GasLimit: 1}, // Wrong gas limit.
	}

	for _, header := range invalidHeaders {
		loopHeader := header
		ingestionFailedResponse := enclave.insertBlockIntoL1Chain(types.NewBlock(&loopHeader, nil, nil, nil, &trie.StackTrie{}))
		if ingestionFailedResponse == nil {
			t.Errorf("expected block with invalid header to be rejected but was accepted")
		}
	}
}
