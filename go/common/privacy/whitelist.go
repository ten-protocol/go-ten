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
	whitelistMap[toEip1967HashHex("eip1967.proxy.implementation")] = true
	whitelistMap[toEip1967FallbackHashHex("org.zeppelinos.proxy.implementation")] = true

	return &Whitelist{
		AllowedStorageSlots: whitelistMap,
	}
}

func toEip1967HashHex(key string) string {
	hash := crypto.Keccak256Hash([]byte(key))
	hashAsBig := hash.Big()
	eipHashHex := "0x" + hashAsBig.Sub(hashAsBig, big.NewInt(1)).Text(16)

	return eipHashHex
}

func toEip1967FallbackHashHex(key string) string {
	hash := crypto.Keccak256Hash([]byte(key))
	hashAsBig := hash.Big()
	eipHashHex := "0x" + hashAsBig.Text(16)

	return eipHashHex
}
