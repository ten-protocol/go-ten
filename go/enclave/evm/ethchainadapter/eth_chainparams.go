package ethchainadapter

import (
	"math"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// ChainParams defines the forks of the EVM machine
// TEN should typically be on the last fork version
func ChainParams(obscuroChainID *big.Int) *params.ChainConfig {
	zeroTimestamp := uint64(0)
	// the forks with this timestamp are not enabled because the current time is always < MaxUint64
	maxTimestamp := uint64(math.MaxUint64)

	// Initialise the database
	return &params.ChainConfig{
		ChainID:             obscuroChainID,
		HomesteadBlock:      gethcommon.Big0,
		DAOForkBlock:        gethcommon.Big0,
		EIP150Block:         gethcommon.Big0,
		EIP155Block:         gethcommon.Big0,
		EIP158Block:         gethcommon.Big0,
		ByzantiumBlock:      gethcommon.Big0,
		ConstantinopleBlock: gethcommon.Big0,
		PetersburgBlock:     gethcommon.Big0,
		IstanbulBlock:       gethcommon.Big0,
		MuirGlacierBlock:    gethcommon.Big0,
		BerlinBlock:         gethcommon.Big0,
		LondonBlock:         gethcommon.Big0,

		CancunTime:   &zeroTimestamp,
		ShanghaiTime: &zeroTimestamp,
		PragueTime:   &zeroTimestamp,
		VerkleTime:   &maxTimestamp, // todo VERKLE - zeroTimestamp,
		OsakaTime:    &maxTimestamp, // todo VERKLE - zeroTimestamp,
	}
}
