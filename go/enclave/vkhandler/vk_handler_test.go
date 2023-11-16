package vkhandler

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
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

	tests := map[string]string{
		"WEMessageFormatTest": viewingkey.GenerateSignMessage(vkPubKeyBytes),
		"OGMessageFormatTest": viewingkey.GenerateSignMessageOG(vkPubKeyBytes, &userAddr),
	}

	for testName, msgToSign := range tests {
		t.Run(testName, func(t *testing.T) {
			signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), userPrivKey)
			assert.NoError(t, err)

			// Create a new vk Handler
			_, err = New(&userAddr, vkPubKeyBytes, signature)
			assert.NoError(t, err)
		})
	}
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
	vkPubKeyBytes := crypto.CompressPubkey(ecies.ImportECDSAPublic(&vkPrivKey.PublicKey).ExportECDSA())

	tests := map[string]string{
		"WEMessageFormatTest": viewingkey.GenerateSignMessage(vkPubKeyBytes),
		"OGMessageFormatTest": viewingkey.GenerateSignMessageOG(vkPubKeyBytes, &userAddr),
	}

	for testName, msgToSign := range tests {
		t.Run(testName, func(t *testing.T) {
			// sign the message
			signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), userPrivKey)
			assert.NoError(t, err)

			// Recover the key based on the signed message and the signature.
			recoveredAccountPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), signature)
			assert.NoError(t, err)
			recoveredAccountAddress := crypto.PubkeyToAddress(*recoveredAccountPublicKey)

			if recoveredAccountAddress.Hex() != userAddr.Hex() {
				t.Errorf("unable to recover user address from signature")
			}

			_, err = crypto.DecompressPubkey(vkPubKeyBytes)
			assert.NoError(t, err)
		})
	}
}
