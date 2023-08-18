package crypto

import (
	"bytes"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/compression"
	"testing"

	gethlog "github.com/ethereum/go-ethereum/log"
)

func TestEncryptionDecryption(t *testing.T) {
	mockData := []byte("mock data for testing")

	logger := gethlog.New("module", "crypto/test")
	service := NewDataEncryptionService(logger)

	// Encrypt the mock data.
	encryptedData := service.Encrypt(mockData)
	if len(encryptedData) == 0 {
		t.Fatal("Failed to encrypt data.")
	}

	testEncryptedData, err := base64.StdEncoding.DecodeString("Lwitrl2FlSymeVlzXCk7JAJnQk3kRtKayF5Of2qtOOhISDtBMFBwY2okYMYbOK/8GajcD+xWqVqm9ztvx5EjYV27FmUqQEjYXoCClMlbqC6AhUtMLYjwt9TPUE02x08CAq2VPAmja8xOAr6WaC8DAyX2s+dC1Fj/gpZ0ZAY7svNdFZjYBw==")
	if err != nil {
		t.Fatal(err)
	}

	// Decrypt the encrypted data.
	decryptedData := service.Decrypt(testEncryptedData)
	if len(decryptedData) == 0 {
		t.Fatal("Failed to decrypt data.")
	}

	comp := compression.NewBrotliDataCompressionService()

	encoded, err := comp.Decompress(decryptedData)
	if err != nil {
		t.Fatal(err)
	}
	var txs []*common.L2Tx
	err = rlp.DecodeBytes(encoded, &txs)
	if err != nil {
		t.Fatal(err)
	}

	// Check if decrypted data matches original mock data.
	if !bytes.Equal(decryptedData, mockData) {
		t.Fatalf("Decrypted data does not match original. Expected %s but got %s", string(mockData), string(decryptedData))
	}
}
