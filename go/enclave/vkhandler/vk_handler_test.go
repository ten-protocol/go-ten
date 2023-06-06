package vkhandler

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/stretchr/testify/assert"
)

func TestVKHandler(t *testing.T) {
	// generate user private Key
	userPrivKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	userAddr := crypto.PubkeyToAddress(userPrivKey.PublicKey)

	// generate ViewingKey private Key
	vkPrivKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	vkPubKeyBytes := crypto.CompressPubkey(ecies.ImportECDSAPublic(&vkPrivKey.PublicKey).ExportECDSA())

	// Sign key
	signature, err := viewingkey.Sign(userPrivKey, vkPubKeyBytes)
	assert.NoError(t, err)

	// Create a new vk Handler
	_, err = New(&userAddr, vkPubKeyBytes, signature)
	assert.NoError(t, err)
}

func TestSignAndCheckSignature(t *testing.T) {
	// generate user private Key
	userPrivKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	userAddr := crypto.PubkeyToAddress(userPrivKey.PublicKey)

	// generate ViewingKey private Key
	vkPrivKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	vkPubKeyByte := crypto.CompressPubkey(ecies.ImportECDSAPublic(&vkPrivKey.PublicKey).ExportECDSA())

	// sign the message
	msgToSign := viewingkey.SignedMsgPrefix + string(vkPubKeyByte)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), userPrivKey)
	assert.NoError(t, err)

	// Recover the key based on the signed message and the signature.
	recoveredAccountPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), signature)
	assert.NoError(t, err)
	recoveredAccountAddress := crypto.PubkeyToAddress(*recoveredAccountPublicKey)

	if recoveredAccountAddress.Hex() != userAddr.Hex() {
		t.Errorf("unable to recover user address from signature")
	}

	_, err = crypto.DecompressPubkey(vkPubKeyByte)
	assert.NoError(t, err)
}
