package walletmock

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type Wallet struct {
	Address common.Address
	// TODO - Store key securely. Geth stores the key encrypted on disk.
	Key keystore.Key
}

func New() Wallet {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(fmt.Sprintf("Could not generate keypair for wallet: %v", err))
	}

	key := keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}

	return Wallet{Address: key.Address, Key: key}
}
