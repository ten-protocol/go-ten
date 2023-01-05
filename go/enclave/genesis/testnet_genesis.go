package genesis

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const TestnetPrefundedPK = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The genesis main account private key.

var TestnetGenesis = Genesis{
	Accounts: []Account{
		{
			Address: gethcommon.HexToAddress("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77"),
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
