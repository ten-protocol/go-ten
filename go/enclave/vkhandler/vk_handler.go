package vkhandler

import (
	"crypto/rand"
	"fmt"
	"io"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"gitlab.com/NebulousLabs/fastrand"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

// AuthenticatedViewingKey - the enclave side of the viewing key. Used for authenticating requests and for encryption
type AuthenticatedViewingKey struct {
	rpcVK          *viewingkey.RPCSignedViewingKey
	AccountAddress *gethcommon.Address
	ecdsaKey       *ecies.PublicKey
	UserID         string
}

func VerifyViewingKey(rpcVK *viewingkey.RPCSignedViewingKey, chainID int64) (*AuthenticatedViewingKey, error) {
	vkPubKey, err := crypto.DecompressPubkey(rpcVK.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decompress viewing key bytes - %w", err)
	}

	rvk := &AuthenticatedViewingKey{
		rpcVK:    rpcVK,
		ecdsaKey: ecies.ImportECDSAPublic(vkPubKey),
	}

	// 2. Authenticate
	recoveredAccountAddress, err := checkViewingKeyAndRecoverAddress(rvk, chainID)
	if err != nil {
		return nil, err
	}

	rvk.AccountAddress = recoveredAccountAddress
	return rvk, nil
}

// checkViewingKeyAndRecoverAddress checks the signature and recovers the address from the viewing key
func checkViewingKeyAndRecoverAddress(vk *AuthenticatedViewingKey, chainID int64) (*gethcommon.Address, error) {
	// get userID from viewingKey public key
	userID := viewingkey.CalculateUserIDHex(vk.rpcVK.PublicKey)
	vk.UserID = userID

	// todo - remove this when the legacy format is no longer supported
	// this is a temporary fix to support the legacy format which will be removed soon
	if vk.rpcVK.SignatureType == viewingkey.Legacy {
		userID = string(vk.rpcVK.PublicKey) // for legacy format, the userID is the public key
	}

	// check the signature and recover the address assuming the message was signed with EIP712
	recoveredSignerAddress, err := viewingkey.CheckSignature(userID, vk.rpcVK.SignatureWithAccountKey, chainID, vk.rpcVK.SignatureType)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed %w", err)
	}

	return recoveredSignerAddress, err
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
