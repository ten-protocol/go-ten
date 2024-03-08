package vkhandler

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
)

const chainID = 443

func TestVKHandler(t *testing.T) {
	// generate user private Key
	userPrivKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	userAccAddress := crypto.PubkeyToAddress(userPrivKey.PublicKey)

	// generate ViewingKey private Key
	vkPrivKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	vkPubKeyBytes := crypto.CompressPubkey(ecies.ImportECDSAPublic(&vkPrivKey.PublicKey).ExportECDSA())
	userID := viewingkey.CalculateUserIDHex(vkPubKeyBytes)
	WEMessageFormatTestHash := accounts.TextHash([]byte(viewingkey.GenerateSignMessage(vkPubKeyBytes)))
	EIP712MessageDataOptions, err := viewingkey.GenerateAuthenticationEIP712RawDataOptions(userID, chainID)
	PersonalSignMessage := viewingkey.GeneratePersonalSignMessage(userID, chainID, 1)
	PersonalSignMessageHash := accounts.TextHash([]byte(PersonalSignMessage))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(EIP712MessageDataOptions) == 0 {
		t.Fatalf("GenerateAuthenticationEIP712RawDataOptions returned no results")
	}
	EIP712MessageFormatTestHash := crypto.Keccak256(EIP712MessageDataOptions[0])

	tests := map[string][]byte{
		"WEMessageFormatTest":           WEMessageFormatTestHash,
		"EIP712MessageFormatTest":       EIP712MessageFormatTestHash,
		"PersonalSignMessageFormatTest": PersonalSignMessageHash,
	}

	for testName, msgHashToSign := range tests {
		t.Run(testName, func(t *testing.T) {
			signature, err := crypto.Sign(msgHashToSign, userPrivKey)
			assert.NoError(t, err)

			// Create a new vk Handler
			_, err = VerifyViewingKey(&viewingkey.RPCSignedViewingKey{
				Account:                 &userAccAddress,
				PublicKey:               vkPubKeyBytes,
				SignatureWithAccountKey: signature,
				SignatureType:           viewingkey.EIP712Signature, // todo - fix this test
			}, chainID)
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
	userID := viewingkey.CalculateUserIDHex(vkPubKeyBytes)
	WEMessageFormatTestHash := accounts.TextHash([]byte(viewingkey.GenerateSignMessage(vkPubKeyBytes)))
	EIP712MessageData, err := viewingkey.GenerateAuthenticationEIP712RawDataOptions(userID, chainID)
	if err != nil {
		t.Fatalf(err.Error())
	}
	EIP712MessageFormatTestHash := crypto.Keccak256(EIP712MessageData[0])
	PersonalSignMessage := viewingkey.GeneratePersonalSignMessage(userID, chainID, 1)
	PersonalSignMessageHash := accounts.TextHash([]byte(PersonalSignMessage))

	tests := map[string][]byte{
		"WEMessageFormatTest":           WEMessageFormatTestHash,
		"EIP712MessageFormatTest":       EIP712MessageFormatTestHash,
		"PersonalSignMessageFormatTest": PersonalSignMessageHash,
	}

	for testName, msgHashToSign := range tests {
		t.Run(testName, func(t *testing.T) {
			// sign the message
			signature, err := crypto.Sign(msgHashToSign, userPrivKey)
			assert.NoError(t, err)

			// Recover the key based on the signed message and the signature.
			recoveredAccountPublicKey, err := crypto.SigToPub(msgHashToSign, signature)
			assert.NoError(t, err)
			recoveredAccountAddress := crypto.PubkeyToAddress(*recoveredAccountPublicKey)

			if recoveredAccountAddress.Hex() != userAddr.Hex() {
				t.Fatalf("Expected user address %s, got %s", userAddr.Hex(), recoveredAccountAddress.Hex())
			}

			_, err = crypto.DecompressPubkey(vkPubKeyBytes)
			assert.NoError(t, err)
		})
	}
}
