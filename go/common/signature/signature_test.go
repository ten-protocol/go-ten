package signature

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
)

// test an address can be recovered from a signature
func TestRecoverAddressFromSignature(t *testing.T) {
	// create a message
	message := []byte("hello world")

	pkStr := hex.EncodeToString(datagenerator.RandomBytes(32))
	privateKey, err := crypto.HexToECDSA(pkStr)
	if err != nil {
		t.Fatal(err)
	}
	messageHash := crypto.Keccak256Hash(message).Bytes()
	// sign the hash
	sig, err := Sign(messageHash, privateKey)
	if err != nil {
		panic(err)
	}

	// recover the address from the signature
	recoveredAddress, err := RecoverAddress(messageHash, sig)
	if err != nil {
		panic(err)
	}

	// check the address is the same as the public key
	if *recoveredAddress != crypto.PubkeyToAddress(privateKey.PublicKey) {
		t.Fatal("address does not match public key")
	}
}
