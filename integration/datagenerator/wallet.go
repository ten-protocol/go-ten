package datagenerator

import (
	"encoding/hex"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/wallet"
)

// RandomWallet returns a wallet with a random private key
func RandomWallet(chainID int64) wallet.Wallet {
	pk := randomHex(32)
	walletConfig := config.HostConfig{
		PrivateKeyString: pk,
		L1ChainID:        chainID,
	}
	return wallet.NewInMemoryWalletFromConfig(walletConfig)
}

func randomHex(n int) string {
	return hex.EncodeToString(RandomBytes(n))
}
