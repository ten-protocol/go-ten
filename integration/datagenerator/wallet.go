package datagenerator

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
)

// RandomWallet returns a wallet with a random private key
func RandomWallet() wallet.Wallet {
	pk, err := randomHex(32)
	if err != nil {
		panic(err) // this should never panic - world should stop if it does
	}
	return wallet.NewInMemoryWallet(pk)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
