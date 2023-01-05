package genesis

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const TestnetPrefundedPK = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The genesis main account private key.

var TestnetGenesis = Genesis{
	Accounts: []Account{
		{
			Address: gethcommon.HexToAddress("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77"),
			Amount:  big.NewInt(7_500_000_000_000_000_000),
		},
	},
}
