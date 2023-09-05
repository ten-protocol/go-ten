package genesis

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const TestnetPrefundedPK = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The genesis main account private key.
var GasBridgingKeys, _ = crypto.GenerateKey()                                                 // todo - make static

var TestnetGenesis = Genesis{
	Accounts: []Account{
		{
			Address: gethcommon.HexToAddress("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77"),
			Amount:  parseHugeNumber("7500000000000000000000000000000"),
		},
		// todo (@stefan) - remove the following when the bridge is updated!
		{ // Address for HOC owner
			Address: gethcommon.HexToAddress("0x987E0a0692475bCc5F13D97E700bb43c1913EFfe"),
			Amount:  parseHugeNumber("7500000000000000000000000000000"),
		},
		{ // Address for POC owner
			Address: gethcommon.HexToAddress("0xDEe530E22045939e6f6a0A593F829e35A140D3F1"),
			Amount:  parseHugeNumber("7500000000000000000000000000000"),
		},
	},
}

// parseHugeNumber parses number that overflow int64
func parseHugeNumber(number string) *big.Int {
	numb, ok := big.NewInt(0).SetString(number, 10)
	if !ok {
		panic(fmt.Sprintf("unable to parse number %s", number))
	}
	return numb
}
