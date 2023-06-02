package datagenerator

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	enclaverpc "github.com/obscuronet/go-obscuro/go/enclave/rpc"
)

func RandomViewingKey() (*viewingkey.ViewingKey, error) {
	keyPair, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key - %w", err)
	}
	pubKey := crypto.CompressPubkey(&keyPair.PublicKey)

	address := crypto.PubkeyToAddress(keyPair.PublicKey)
	return &viewingkey.ViewingKey{
		Account:    &address,
		PrivateKey: ecies.ImportECDSA(keyPair),
		PublicKey:  pubKey,
		SignedKey:  viewingkey.Sign(), // this guys needs to know the priv side that the WE generated
	}, err
}

// Signs a viewing key.
func signViewingKey(privateKey *ecdsa.PrivateKey, viewingKey []byte) []byte {
	msgToSign := enclaverpc.ViewingKeySignedMsgPrefix + string(viewingKey)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), privateKey)
	if err != nil {
		panic(err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	return signatureWithLeadBytes
}
