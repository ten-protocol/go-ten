package vkhandler

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/accounts"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"gitlab.com/NebulousLabs/fastrand"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

// AuthenticatedViewingKey - This data accompanies every RPC request.
type AuthenticatedViewingKey struct {
	VkPubKey       []byte
	Signature      []byte
	AccountAddress *gethcommon.Address
	ecdsaKey       *ecies.PublicKey
}

// ExtractAndAuthenticateViewingKey extracts the VK from the request and authenticates the signature against the account key
func ExtractAndAuthenticateViewingKey(rawVk interface{}, chainID int64) (*AuthenticatedViewingKey, error) {
	// 1. Extract
	vkBytesList, ok := rawVk.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to cast the vk to []any")
	}

	if len(vkBytesList) != 2 {
		return nil, fmt.Errorf("wrong size of viewing key params")
	}

	vkPubkey, err := hexutil.Decode(vkBytesList[0].(string))
	if err != nil {
		return nil, fmt.Errorf("could not decode data in vk pub key - %w", err)
	}

	accountSignatureHex, err := hexutil.Decode(vkBytesList[1].(string))
	if err != nil {
		return nil, fmt.Errorf("could not decode data in vk signature - %w", err)
	}

	// 2. Authenticate
	return AuthenticateViewingKey(vkPubkey, accountSignatureHex, chainID)
}

func AuthenticateViewingKey(vkPubkey []byte, accountSignatureHex []byte, chainID int64) (*AuthenticatedViewingKey, error) {
	vkPubKey, err := crypto.DecompressPubkey(vkPubkey)
	if err != nil {
		return nil, fmt.Errorf("could not decompress viewing key bytes - %w", err)
	}

	rvk := &AuthenticatedViewingKey{
		VkPubKey:  vkPubkey,
		ecdsaKey:  ecies.ImportECDSAPublic(vkPubKey),
		Signature: accountSignatureHex,
	}

	// 2. Authenticate
	recoveredAccountAddress, err := checkViewingKeyAndRecoverAddress(rvk, chainID)
	if err != nil {
		return nil, err
	}

	rvk.AccountAddress = recoveredAccountAddress
	return rvk, nil
}

func checkViewingKeyAndRecoverAddress(vk *AuthenticatedViewingKey, chainID int64) (*gethcommon.Address, error) {
	// get userID from viewingKey public key
	userID := viewingkey.CalculateUserIDHex(vk.VkPubKey)

	// check signature and recover the address
	address, err := viewingkey.CheckEIP712Signature(userID, vk.Signature, chainID) //nolint:ineffassign
	if err != nil {
		// try the legacy format
		legacyMessage := viewingkey.GenerateSignMessage(vk.VkPubKey)
		legacyMessageHash := accounts.TextHash([]byte(legacyMessage))
		address, err = viewingkey.CheckSignatureAndReturnAccountAddress(legacyMessageHash, vk.Signature)
		if err == nil {
			return address, nil
		}
	}

	// TODO @Ziga - this must be removed.
	msgToSignLegacy := viewingkey.GenerateSignMessage(vk.VkPubKey)
	recoveredAccountPublicKeyLegacy, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSignLegacy)), vk.Signature)
	if err != nil {
		return nil, fmt.Errorf("viewing key but could not validate its signature - %w", err)
	}
	recoveredAccountAddressLegacy := crypto.PubkeyToAddress(*recoveredAccountPublicKeyLegacy)
	address = &recoveredAccountAddressLegacy

	return address, err
}

// crypto.rand is quite slow. When this variable is true, we will use a fast CSPRNG algorithm
const useFastRand = true

func rndSource() io.Reader {
	rndSource := rand.Reader
	if useFastRand {
		rndSource = fastrand.Reader
	}
	return rndSource
}

// Encrypt returns the payload encrypted with the viewingKey
func (vk *AuthenticatedViewingKey) Encrypt(bytes []byte) ([]byte, error) {
	if len(bytes) == 0 {
		bytes = placeholderResult
	}
	encryptedBytes, err := ecies.Encrypt(rndSource(), vk.ecdsaKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt with given public VK - %w", err)
	}

	return encryptedBytes, nil
}
