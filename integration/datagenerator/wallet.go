package datagenerator

import (
	"encoding/hex"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/wallet"
)

// RandomWallet returns a wallet with a random private key
func RandomWallet(chainID int64) wallet.Wallet {
	return wallet.NewInMemoryWalletFromConfig(
		randomHex(32),
		chainID,
		testlog.Logger(),
	)
}

func randomHex(n int) string {
	return hex.EncodeToString(RandomBytes(n))
}
