package components

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

const (
	InitialBaseFee        = params.InitialBaseFee / 10 // InitialBaseFee is ETH's base fee.
	ComponentName         = "gas_pricing"
	DynamicPricingTrigger = "dynamic-pricing"
	StaticPricingTrigger  = "static_pricing"
)

type GasPricer struct {
	logger  gethlog.Logger
	config  *config.EnclaveConfig
	storage storage.Storage
}

func NewGasPricer(logger gethlog.Logger, config *config.EnclaveConfig, storage storage.Storage) *GasPricer {
	return &GasPricer{
		logger:  logger,
		config:  config,
		storage: storage,
	}
}

// CalculateBlockBaseFeeAtHeight calculates the base fee for a block using height-based pricing logic
func (gp *GasPricer) CalculateBlockBaseFeeAtHeight(ctx context.Context, cfg *params.ChainConfig, parent *types.Header, l1Height uint64) *big.Int {
	// Dynamic pricing is always enabled - use the EIP-1559 calculation
	if parent == nil {
		return new(big.Int).SetUint64(InitialBaseFee)
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

func (gp *GasPricer) StaticL2BaseFee(header *types.Header) *big.Int {
	return gp.config.BaseFee
}
