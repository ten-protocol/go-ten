package components

import (
	"context"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/enclave/config"
)

const (
	InitialBaseFee        = params.InitialBaseFee / 10 // InitialBaseFee is ETH's base fee.
	ComponentName         = "gas_pricing"
	DynamicPricingTrigger = "gas_pricing"
	StaticPricingTrigger  = "static_pricing"
)

type GasPricer struct {
	logger                gethlog.Logger
	config                *config.EnclaveConfig
	dynamicPricingEnabled atomic.Bool // tracks whether dynamic pricing upgrade has been applied
}

func NewGasPricer(logger gethlog.Logger, config *config.EnclaveConfig) *GasPricer {
	return &GasPricer{
		logger:                logger,
		config:                config,
		dynamicPricingEnabled: atomic.Bool{}, // start with static pricing (minGasPrice)
	}
}

// CanUpgrade implements the UpgradeHandler interface for gas pricing upgrades
func (gp *GasPricer) CanUpgrade(ctx context.Context, featureName string, featureData []byte) bool {
	// Gas pricer only handles gas_pricing feature upgrades
	if featureName != "gas_pricing" {
		return false
	}

	// Check if we support the specific upgrade data
	switch string(featureData) {
	case "dynamic-pricing":
		return true
	default:
		gp.logger.Warn("Unsupported gas pricing upgrade data", "featureData", string(featureData))
		return false
	}
}

// HandleUpgrade implements the UpgradeHandler interface for gas pricing upgrades
func (gp *GasPricer) HandleUpgrade(ctx context.Context, featureName string, featureData []byte) error {
	gp.logger.Info("Processing gas pricing upgrade",
		"featureName", featureName,
		"featureData", string(featureData))

	// Check if this is the dynamic pricing upgrade
	if string(featureData) == DynamicPricingTrigger {
		gp.dynamicPricingEnabled.Store(true)
		gp.logger.Info("Dynamic gas pricing enabled")
	} else if string(featureData) == StaticPricingTrigger {
		gp.dynamicPricingEnabled.Store(false)
		gp.logger.Info("Static gas pricing enabled")
	} else {
		gp.logger.Warn("Unknown gas pricing upgrade data", "featureData", string(featureData))
	}

	return nil
}

func (gp *GasPricer) CalculateBlockBaseFee(cfg *params.ChainConfig, parent *types.Header) *big.Int {
	// If dynamic pricing is not enabled, return the configured minimum gas price
	if !gp.dynamicPricingEnabled.Load() {
		gp.logger.Trace("Using static gas pricing", "minGasPrice", gp.config.MinGasPrice)
		return new(big.Int).Set(gp.config.MinGasPrice)
	}

	// Dynamic pricing is enabled - use the EIP-1559 calculation
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
