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

// GenerateAndRegisterViewingKey take an obscuro client and a wallet, it generates a keypair, simulates signing the key,
//	configures it on the client and registers the viewing key with the enclave
func GenerateAndRegisterViewingKey(cli rpcclientlib.Client, wal wallet.Wallet) {
	vk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	viewingPrivateKeyEcies := ecies.ImportECDSA(vk)
	viewingPubKey := crypto.CompressPubkey(&vk.PublicKey)
	cli.SetViewingKey(viewingPrivateKeyEcies, viewingPubKey)

	viewingKeyHex := hex.EncodeToString(viewingPubKey)

	signature := signViewingKey(viewingKeyHex, wal.PrivateKey())

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
