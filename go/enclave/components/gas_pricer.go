package components

import (
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/enclave/config"
)

const (
	InitialBaseFee = params.InitialBaseFee // modify this to something sensible
)

type GasPricer struct {
	logger gethlog.Logger
	config *config.EnclaveConfig
}

func NewGasPricer(logger gethlog.Logger, config *config.EnclaveConfig) *GasPricer {
	return &GasPricer{
		logger: logger,
		config: config,
	}
}

func (gp *GasPricer) CalculateBlockBaseFee(cfg *params.ChainConfig, parent *types.Header) *big.Int {
	if parent == nil {
		return big.NewInt(InitialBaseFee)
	}

	// This uses the parent gas limit, divided by cfg.ElasticityMultiplier()
	// to determine the gas target.
	// If the target is equal to used baseFee is unchanged.
	// Else its moved up and down using cfg.BaseFeeChangeDenominator()
	// to determine the amount of change.
	calculatedBaseFee := eip1559.CalcBaseFee(cfg, parent)

	// Ensure the base fee never falls below the configured minimum
	if calculatedBaseFee.Cmp(gp.config.MinGasPrice) < 0 {
		return new(big.Int).Set(gp.config.MinGasPrice)
	}

	return calculatedBaseFee
}
