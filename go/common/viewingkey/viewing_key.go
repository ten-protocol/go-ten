package viewingkey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
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

// GenerateViewingKeyForWallet takes an account wallet, generates a viewing key and signs the key with the acc's private key
// uses the same method of signature handling as Metamask/geth
// TODO @Ziga - update this method to use the new EIP-712 signature format / personal sign after the removal of the legacy format
func GenerateViewingKeyForWallet(wal wallet.Wallet) (*ViewingKey, error) {
	// generate an ECDSA key pair to encrypt sensitive communications with the obscuro enclave
	vk, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate viewing key for RPC client: %w", err)
	}

	// get key in ECIES format
	viewingPrivateKeyECIES := ecies.ImportECDSA(vk)

	// encode public key as bytes
	viewingPubKeyBytes := crypto.CompressPubkey(&vk.PublicKey)

	// sign public key bytes with the wallet's private key
	signature, err := mmSignViewingKey(viewingPubKeyBytes, wal.PrivateKey())
	if err != nil {
		return nil, err
	}

	accAddress := wal.Address()
	return &ViewingKey{
		Account:                 &accAddress,
		PrivateKey:              viewingPrivateKeyECIES,
		PublicKey:               viewingPubKeyBytes,
		SignatureWithAccountKey: signature,
		SignatureType:           Legacy,
	}, nil
}

// mmSignViewingKey takes a public key bytes as hex and the private key for a wallet, it simulates the back-and-forth to
// MetaMask and returns the signature bytes to register with the enclave
func mmSignViewingKey(viewingPubKeyBytes []byte, signerKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := Sign(signerKey, viewingPubKeyBytes)
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

// Sign takes a users Private key and signs the public viewingKey hex
func Sign(userPrivKey *ecdsa.PrivateKey, vkPubKey []byte) ([]byte, error) {
	msgToSign := GenerateSignMessage(vkPubKey)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign messages - %w", err)
	}
	return signature, nil
}
