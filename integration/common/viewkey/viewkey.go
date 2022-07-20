package viewkey

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/enclave/rpcencryptionmanager"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// This package contains viewing key utils for testing and simulations

// GenerateAndRegisterViewingKey takes a wallet and the client that will be used for that wallet's transactions.
//
//	It sets up a viewing key for that client (without using a wallet extension) with the following steps:
//	1. generate a "viewing" keypair
//	2. simulates signing the key with metamask
//	3. sets the private key on the client (to decrypt enclave responses)
//	4. registers the public viewing key with the enclave (to encrypt enclave responses)
func GenerateAndRegisterViewingKey(cli rpcclientlib.Client, wal wallet.Wallet) {
	// generate an ECDSA key pair to encrypt sensitive communications with the obscuro enclave
	vk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	// get key in ECIES format
	viewingPrivateKeyECIES := ecies.ImportECDSA(vk)

	// encode public key as bytes
	viewingPubKeyBytes := crypto.CompressPubkey(&vk.PublicKey)

	// set key pair on the RPC client
	cli.SetViewingKey(viewingPrivateKeyECIES, viewingPubKeyBytes)

	// sign hex-encoded public key string with the wallet's private key
	viewingKeyHex := hex.EncodeToString(viewingPubKeyBytes)
	signature := signViewingKey(viewingKeyHex, wal.PrivateKey())

	// submit the signed public key to the enclave so it can encrypt sensitive responses
	err = cli.RegisterViewingKey(wal.Address(), signature)
	if err != nil {
		panic(err)
	}
}

// signViewingKey takes a public key bytes as hex and the private key for a wallet, it simulates the back-and-forth to
//	metamask and returns the signature bytes to register with the enclave
func signViewingKey(viewingKeyHex string, signerKey *ecdsa.PrivateKey) []byte {
	msgToSign := rpcencryptionmanager.ViewingKeySignedMsgPrefix + viewingKeyHex
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), signerKey)
	if err != nil {
		panic(err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	// this string encoded signature is what the wallet extension would receive after it is signed by metamask
	sigStr := hex.EncodeToString(signatureWithLeadBytes)
	// and then we extract the signature bytes in the same way as the wallet extension
	outputSig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		panic(err)
	}
	// This same change is made in geth internals, for legacy reasons to be able to recover the address:
	//	https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	outputSig[64] -= 27

	return outputSig
}
