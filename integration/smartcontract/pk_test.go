package smartcontract

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

func TestSign(t *testing.T) {
	privateKeyA, err := crypto.ToECDSA(datagenerator.RandomBytes(32))
	if err != nil {
		t.Error(err)
	}
	pubKeyA := privateKeyA.PublicKey
	addrA := crypto.PubkeyToAddress(pubKeyA)

	// pk signs a random message
	msg := crypto.Keccak256([]byte("foo"))
	sig, err := crypto.Sign(msg, privateKeyA)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}

	// the pubkey recovered is the same as the pk.pubkey ?
	recoveredPub, err := crypto.Ecrecover(msg, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	recoveredPubKey, err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		t.Error(err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)
	if !bytes.Equal(addrA.Bytes(), recoveredAddr.Bytes()) {
		t.Errorf("unexpected addresses Got %s, expected: %s", addrA.Hex(), recoveredAddr.Hex())
	}

	ecrecover, err := crypto.Ecrecover(msg, sig)
	if err != nil {
		t.Error(err)
	}
	recoveredPubKey, err = crypto.UnmarshalPubkey(ecrecover)
	if err != nil {
		t.Error(err)
	}
	recoveredAddr = crypto.PubkeyToAddress(*recoveredPubKey)

	if !bytes.Equal(addrA.Bytes(), recoveredAddr.Bytes()) {
		t.Errorf("unexpected addresses Got %s, expected: %s", addrA.Hex(), recoveredAddr.Hex())
	}
}

func TestRetrieve(t *testing.T) {
	privateKeyA, err := crypto.ToECDSA(datagenerator.RandomBytes(32))
	if err != nil {
		t.Error(err)
	}
	pubKeyA := privateKeyA.PublicKey
	addrA := crypto.PubkeyToAddress(pubKeyA)
	t.Logf("pubkey addr: %s", addrA.Hex())
	pubKeyABytes := crypto.FromECDSAPub(&pubKeyA)

	// pk signs a random message
	msg := crypto.Keccak256([]byte("foo"))
	sig, err := crypto.Sign(msg, privateKeyA)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}

	// recover the pubkey given the msg and the signature
	recoveredPub, err := crypto.Ecrecover(msg, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	recoveredPubKey, err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(pubKeyABytes, recoveredPub) {
		t.Errorf("unexpected pub keys. got: %s, expected: %s", recoveredPub, pubKeyABytes)
	}
	// recovered Pubkey yeilds the same address
	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)
	if !bytes.Equal(addrA.Bytes(), recoveredAddr.Bytes()) {
		t.Errorf("unexpected addresses Got %s, expected: %s", addrA.Hex(), recoveredAddr.Hex())
	}

	sig = sig[:len(sig)-1] // remove recovery id
	verified := crypto.VerifySignature(pubKeyABytes, msg, sig)
	if !verified {
		t.Error("not verified")
	}
}
