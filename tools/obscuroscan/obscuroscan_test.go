package obscuroscan

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/crypto"

	"github.com/obscuronet/obscuro-playground/integration/datagenerator"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func TestCanDecryptTxBlob(t *testing.T) {
	txs := []*nodecommon.L2Tx{datagenerator.CreateL2Tx(), datagenerator.CreateL2Tx()}

	txsJSONBytes, err := decryptTxBlob(generateEncryptedTxBlob(txs))
	if err != nil {
		t.Fatalf("transaction blob decryption failed. Cause: %s", err)
	}

	expectedTxsJSONBytes, err := json.Marshal(txs)
	if err != nil {
		t.Fatalf("marshalling transactions to JSON failed. Cause: %s", err)
	}

	if string(expectedTxsJSONBytes) != string(txsJSONBytes) {
		t.Fatalf("expected %s, got %s", string(expectedTxsJSONBytes), string(txsJSONBytes))
	}
}

func TestThrowsIfEncryptedRollupIsInvalid(t *testing.T) {
	_, err := decryptTxBlob([]byte("invalid_tx_blob"))
	if err == nil {
		t.Fatal("did not error on invalid transaction blob")
	}
}

// Generates an encrypted transaction blob in Base64 encoding.
func generateEncryptedTxBlob(txs []*nodecommon.L2Tx) []byte {
	txBlob := crypto.NewTransactionBlobCryptoImpl().Encrypt(txs)
	return []byte(base64.StdEncoding.EncodeToString(txBlob))
}
