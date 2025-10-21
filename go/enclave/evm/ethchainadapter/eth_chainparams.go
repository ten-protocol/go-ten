package ethchainadapter

import (
	"math"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// ChainParams defines the forks of the EVM machine
// TEN should typically be on the last fork version
func ChainParams(tenChainID *big.Int) *params.ChainConfig {
	zeroTimestamp := uint64(0)

	return &params.ChainConfig{
		ChainID: tenChainID,

		HomesteadBlock: gethcommon.Big0,
		DAOForkBlock:   nil,
		DAOForkSupport: false,

		EIP150Block: gethcommon.Big0,
		EIP155Block: gethcommon.Big0,
		EIP158Block: gethcommon.Big0,

		ByzantiumBlock:      gethcommon.Big0,
		ConstantinopleBlock: gethcommon.Big0,
		PetersburgBlock:     gethcommon.Big0,
		IstanbulBlock:       gethcommon.Big0,
		MuirGlacierBlock:    gethcommon.Big0,
		BerlinBlock:         gethcommon.Big0,
		LondonBlock:         gethcommon.Big0,
		ArrowGlacierBlock:   gethcommon.Big0,
		GrayGlacierBlock:    gethcommon.Big0,
		MergeNetsplitBlock:  nil,

		ShanghaiTime: &zeroTimestamp,
		CancunTime:   &zeroTimestamp,
		PragueTime:   &zeroTimestamp,
		OsakaTime:    &zeroTimestamp,
		VerkleTime:   &zeroTimestamp,

		EnableVerkleAtGenesis:   true,
		TerminalTotalDifficulty: big.NewInt(math.MaxInt64),
	}
}
