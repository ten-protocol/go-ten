package limiters

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func TestGasLimiter(t *testing.T) {
	limiter := NewGasLimiter()

	// block from -  https://etherscan.io/block/17423565
	baseFee := big.NewInt(0)
	gweiFee := big.NewFloat(34.06728554)
	gweiFee.Mul(gweiFee, big.NewFloat(params.GWei))
	gweiFee.Int(baseFee)

	blockHeader := types.Header{
		GasLimit: 30_000_000,
		GasUsed:  14_917_852,
		BaseFee:  baseFee,
	}

	limiter.ProcessBlock(&blockHeader)
	limit := limiter.GetCalldataLimit()
	if limit != 1875000 {
		t.Fail()
	}
}
