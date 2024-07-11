package privacy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type Whitelist struct {
	AllowedStorageSlots map[string]bool
}

func NewWhitelist() *Whitelist {
	whitelistMap := make(map[string]bool)
	whitelistMap[toEip1967HashHex("eip1967.proxy.beacon")] = true

	return &Whitelist{
		AllowedStorageSlots: whitelistMap,
	}
}

func toEip1967HashHex(key string) string {
	hash := crypto.Keccak256Hash([]byte(key))
	hashAsbig := hash.Big()
	eipHashHex := "0x" + hashAsbig.Sub(hashAsbig, big.NewInt(1)).Text(16)

	return eipHashHex
}
