package datagenerator

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ten-protocol/go-ten/go/wallet"
)

// RandomWallet returns a wallet with a random private key
func RandomWallet(chainID int64) wallet.Wallet {
	return wallet.NewInMemoryWalletFromConfig(
		randomHex(32),
		chainID,
		testlog.Logger(),
	)
}

func RandomPrivateKey() *ecdsa.PrivateKey {
	privKey, err := crypto.HexToECDSA(randomHex(32))
	if err != nil {
		panic(err)
	}
	return privKey
}

func randomHex(n int) string {
	return hex.EncodeToString(RandomBytes(n))
}
