package vkhandler

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
)

const chainID = 443

// generateRandomUserKeys -
// generates a random user private key and a random viewing key private key and returns the user private key,
// the viewing key private key, the userID and the user address
func generateRandomUserKeys() (*ecdsa.PrivateKey, *ecdsa.PrivateKey, string, gethcommon.Address) {
	userPrivKey, err := crypto.GenerateKey() // user private key
	if err != nil {
		return nil, nil, "", gethcommon.Address{}
	}
	vkPrivKey, _ := crypto.GenerateKey() // viewingkey generated in the gateway
	if err != nil {
		return nil, nil, "", gethcommon.Address{}
	}

	// get the address from userPrivKey
	userAddress := crypto.PubkeyToAddress(userPrivKey.PublicKey)

	// get userID from viewingKey public key
	vkPubKeyBytes := crypto.CompressPubkey(ecies.ImportECDSAPublic(&vkPrivKey.PublicKey).ExportECDSA())
	userID := viewingkey.CalculateUserIDHex(vkPubKeyBytes)

	return userPrivKey, vkPrivKey, userID, userAddress
}

// MessageWithSignatureType - struct to hold the message hash and the signature type for testing purposes
type MessageWithSignatureType struct {
	Hash          []byte
	SignatureType viewingkey.SignatureType
}

func TestCheckSignature(t *testing.T) {
	userPrivKey, _, userID, userAddress := generateRandomUserKeys()

	// Generate all message types and create map with the corresponding signature type
	// Test EIP712 message format
	EIP712MessageDataOptions, err := viewingkey.GenerateAuthenticationEIP712RawDataOptions(userID, chainID)
	if err != nil {
		t.Fatalf(err.Error())
	}
	EIP712MessageHash := crypto.Keccak256(EIP712MessageDataOptions[0])
	PersonalSignMessageHash := accounts.TextHash([]byte(viewingkey.GeneratePersonalSignMessage(userID, chainID, viewingkey.PersonalSignMessageSupportedVersions[0])))

	messages := map[string]MessageWithSignatureType{
		"EIP712MessageHash": {
			Hash:          EIP712MessageHash,
			SignatureType: viewingkey.EIP712Signature,
		},
		"PersonalSignMessageHash": {
			Hash:          PersonalSignMessageHash,
			SignatureType: viewingkey.PersonalSign,
		},
	}
	// sign each message hash with the user private key and check the signature with the corresponding signature type
	for testName, message := range messages {
		fmt.Println(testName)
		t.Run(testName, func(t *testing.T) {
			signature, err := crypto.Sign(message.Hash, userPrivKey)
			assert.NoError(t, err)

			addr, err := viewingkey.CheckSignature(userID, signature, chainID, message.SignatureType)
			assert.NoError(t, err)

			assert.Equal(t, userAddress.Hex(), addr.Hex())
		})
	}
}

func TestVerifyViewingKey(t *testing.T) {
	userPrivKey, vkPrivKey, userID, userAddress := generateRandomUserKeys()
	fmt.Println("User Address: ", userAddress.Hex())
	// Generate all message types and create map with the corresponding signature type
	// Test EIP712 message format
	EIP712MessageDataOptions, err := viewingkey.GenerateAuthenticationEIP712RawDataOptions(userID, chainID)
	if err != nil {
		t.Fatalf(err.Error())
	}
	EIP712MessageHash := crypto.Keccak256(EIP712MessageDataOptions[0])
	PersonalSignMessageHash := accounts.TextHash([]byte(viewingkey.GeneratePersonalSignMessage(userID, chainID, viewingkey.PersonalSignMessageSupportedVersions[0])))

	messages := map[string]MessageWithSignatureType{
		"EIP712MessageHash": {
			Hash:          EIP712MessageHash,
			SignatureType: viewingkey.EIP712Signature,
		},
		"PersonalSignMessageHash": {
			Hash:          PersonalSignMessageHash,
			SignatureType: viewingkey.PersonalSign,
		},
	}

	for testName, message := range messages {
		t.Run(testName, func(t *testing.T) {
			signature, err := crypto.Sign(message.Hash, userPrivKey)
			assert.NoError(t, err)

			vkPubKeyBytes := crypto.CompressPubkey(ecies.ImportECDSAPublic(&vkPrivKey.PublicKey).ExportECDSA())
			// Create a new vk Handler
			rpcSignedVK, err := VerifyViewingKey(&viewingkey.RPCSignedViewingKey{
				PublicKey:               vkPubKeyBytes,
				SignatureWithAccountKey: signature,
				SignatureType:           message.SignatureType,
			}, chainID)
			assert.NoError(t, err)
			assert.Equal(t, rpcSignedVK.UserID, userID)
			assert.Equal(t, rpcSignedVK.AccountAddress.Hex(), userAddress.Hex())
		})
	}
}
