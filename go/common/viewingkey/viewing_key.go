package viewingkey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/wallet"
)

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave.
// It is the client-side perspective of the viewing key used for decrypting incoming traffic.
type ViewingKey struct {
	Account                 *gethcommon.Address // Account address that this Viewing Key is bound to - Users Pubkey address
	PrivateKey              *ecies.PrivateKey   // ViewingKey private key to encrypt data to the enclave
	PublicKey               []byte              // ViewingKey public key in decrypt data from the enclave
	SignatureWithAccountKey []byte              // ViewingKey public key signed by the Accounts Private key - Allows to retrieve the Account address
	SignatureType           SignatureType       // Type of signature used to sign the public key
}

// RPCSignedViewingKey - used for transporting a minimalist viewing key via
// every RPC request to a sensitive method, including Log subscriptions.
// only the public key and the signature are required
// the account address is sent as well to aid validation
type RPCSignedViewingKey struct {
	PublicKey               []byte
	SignatureWithAccountKey []byte
	SignatureType           SignatureType
}

const (
	pubKeyLen = 33
	sigLen    = 65
)

func (vk RPCSignedViewingKey) Validate() error {
	if len(vk.PublicKey) != pubKeyLen {
		return fmt.Errorf("invalid viewing key")
	}
	if len(vk.SignatureWithAccountKey) != sigLen {
		return fmt.Errorf("invalid viewing key signature")
	}
	return nil
}

// GenerateViewingKeyForWallet takes an account wallet, generates a viewing key and signs the key with the acc's private key
func GenerateViewingKeyForWallet(wal wallet.Wallet) (*ViewingKey, error) {
	chainID := wal.ChainID().Int64()
	messageType := PersonalSign

	// simulate what the gateway would do to generate the viewing key
	viewingKeyPrivate, err := crypto.GenerateKey()
	viewingPrivateKeyECIES := ecies.ImportECDSA(viewingKeyPrivate)
	if err != nil {
		return nil, err
	}
	encryptionToken := CalculateUserID(crypto.CompressPubkey(viewingPrivateKeyECIES.PublicKey.ExportECDSA()))
	messageToSign, err := GenerateMessage(encryptionToken, chainID, PersonalSignVersion, messageType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate message for viewing key: %w", err)
	}
	msgHash, err := GetMessageHash(messageToSign, messageType)
	if err != nil {
		return nil, err
	}

	signature, err := mmSignViewingKey(msgHash, wal.PrivateKey())
	if err != nil {
		return nil, err
	}
	vkPubKeyBytes := crypto.CompressPubkey(ecies.ImportECDSAPublic(&viewingKeyPrivate.PublicKey).ExportECDSA())
	accAddress := wal.Address()
	return &ViewingKey{
		Account:                 &accAddress,
		PrivateKey:              viewingPrivateKeyECIES,
		PublicKey:               vkPubKeyBytes,
		SignatureWithAccountKey: signature,
		SignatureType:           PersonalSign,
	}, nil
}

// mmSignViewingKey takes a public key bytes as hex and the private key for a wallet, it simulates the back-and-forth to
// MetaMask and returns the signature bytes to register with the enclave
func mmSignViewingKey(messageHash []byte, signerKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := crypto.Sign(messageHash, signerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign viewing key: %w", err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	// this string encoded signature is what the wallet extension would receive after it is signed by metamask
	sigStr := hex.EncodeToString(signatureWithLeadBytes)
	// and then we extract the signature bytes in the same way as the wallet extension
	outputSig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature string: %w", err)
	}
	// This same change is made in geth internals, for legacy reasons to be able to recover the address:
	//	https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	outputSig[64] -= 27

	return outputSig, nil
}
