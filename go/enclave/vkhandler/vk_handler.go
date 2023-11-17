package vkhandler

import (
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var ErrInvalidAddressSignature = fmt.Errorf("invalid viewing key signature for requested address")

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

// VKHandler handles encryption and validation of viewing keys
type VKHandler struct {
	publicViewingKey *ecies.PublicKey
}

// VKHandler is responsible for:
// - checking if received signature of a provided viewing key is signed by provided address
// - encrypting payloads with a viewing key (public key) that can only be decrypted by private key signed owned by an address signing it

// New creates a new viewing key handler if signature is valid and was produced by given address
// It receives address, viewing key and a signature over viewing key.
// In order to check signature validity, we need to reproduce a message that was originally signed
func New(requestedAddr *gethcommon.Address, vkPubKeyBytes, accountSignatureHexBytes []byte) (*VKHandler, error) {
	// get userID from viewingKey public key
	userID := viewingkey.CalculateUserIDHex(vkPubKeyBytes)

	// check if a signature is valid and signed address matched requiredAddress
	isValidSignature, _ := viewingkey.VerifySignatureEIP712(userID, requestedAddr, accountSignatureHexBytes)
	// TODO: @ziga
	// don't return here, because we still need to support old Wallet extension endpoints now (will be removed in #2134)
	//if err != nil {
	//	return nil, err
	//}

	// Old wallet extension message format
	// We recover the key based on the signed message and the signature (same as before, but with legacy message format "vk"+<vk>"
	// todo (@ziga) remove support of old message format when removing old wallet extension endpoints (#2134)
	msgToSignLegacy := viewingkey.GenerateSignMessage(vkPubKeyBytes)
	recoveredAccountPublicKeyLegacy, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSignLegacy)), accountSignatureHexBytes)
	if err != nil {
		return nil, fmt.Errorf("viewing key but could not validate its legacy signature - %w", err)
	}
	recoveredAccountAddressLegacy := crypto.PubkeyToAddress(*recoveredAccountPublicKeyLegacy)

	// is the requested account address the same as the address recovered from the signature
	// todo (@ziga) - we currently check also for legacy address and allow both (remove this after transition period)
	if isValidSignature &&
		requestedAddr.Hash() != recoveredAccountAddressLegacy.Hash() {
		return nil, ErrInvalidAddressSignature
	}

	// We decompress the viewing key and create the corresponding ECIES key.
	viewingKey, err := crypto.DecompressPubkey(vkPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not decompress viewing key bytes - %w", err)
	}

	return &VKHandler{
		publicViewingKey: ecies.ImportECDSAPublic(viewingKey),
	}, nil
}

// Encrypt returns the payload encrypted with the viewingKey
func (m *VKHandler) Encrypt(bytes []byte) ([]byte, error) {
	if len(bytes) == 0 {
		bytes = placeholderResult
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, m.publicViewingKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt with given public VK - %w", err)
	}

	return encryptedBytes, nil
}
