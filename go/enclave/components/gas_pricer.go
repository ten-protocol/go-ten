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
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

const (
	InitialBaseFee        = params.InitialBaseFee / 10 // InitialBaseFee is ETH's base fee.
	ComponentName         = "gas_pricing"
	DynamicPricingTrigger = "dynamic-pricing"
	StaticPricingTrigger  = "static_pricing"
)

type GasPricer struct {
	logger                gethlog.Logger
	config                *config.EnclaveConfig
	storage               storage.Storage
	dynamicPricingEnabled atomic.Bool // tracks whether dynamic pricing upgrade has been applied
}

func NewGasPricer(logger gethlog.Logger, config *config.EnclaveConfig, storage storage.Storage) *GasPricer {
	return &GasPricer{
		logger:                logger,
		config:                config,
		storage:               storage,
		dynamicPricingEnabled: atomic.Bool{}, // start with static pricing (minGasPrice)
	}
}

// CanUpgrade implements the UpgradeHandler interface for gas pricing upgrades
func (gp *GasPricer) CanUpgrade(ctx context.Context, featureName string, featureData []byte) bool {
	// Gas pricer only handles gas_pricing feature upgrades
	if featureName != ComponentName {
		return false
	}

	// Check if we support the specific upgrade data
	switch string(featureData) {
	case DynamicPricingTrigger:
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

	// Handle the specific upgrade data
	switch string(featureData) {
	case DynamicPricingTrigger:
		gp.dynamicPricingEnabled.Store(true)
		gp.logger.Info("Dynamic gas pricing enabled")
	case StaticPricingTrigger:
		gp.dynamicPricingEnabled.Store(false)
		gp.logger.Info("Static gas pricing enabled")
	default:
		gp.logger.Warn("Unknown gas pricing upgrade data", "featureData", string(featureData))
	}

	return nil
}

// isDynamicPricingEnabledAtHeight checks if dynamic pricing is enabled for a given L1 block height
func (gp *GasPricer) isDynamicPricingEnabledAtHeight(ctx context.Context, l1Height uint64) bool {
	// Query for finalized dynamic pricing upgrades
	finalizedUpgrades, err := gp.storage.GetFinalizedNetworkUpgrades(ctx)
	if err != nil {
		gp.logger.Warn("Failed to get finalized network upgrades, falling back to atomic flag", "error", err)
		return gp.dynamicPricingEnabled.Load()
	}

	// Check if dynamic pricing upgrade is activated at this height
	for _, upgrade := range finalizedUpgrades {
		if upgrade.FeatureName == ComponentName && string(upgrade.FeatureData) == DynamicPricingTrigger {
			// The upgrade is activated if the L1 height is >= the finalized activation height
			if l1Height >= upgrade.AppliedAtL1Height {
				return true
			}
		}
	}

	return false
}

// CalculateBlockBaseFeeAtHeight calculates the base fee for a block using height-based pricing logic
func (gp *GasPricer) CalculateBlockBaseFeeAtHeight(ctx context.Context, cfg *params.ChainConfig, parent *types.Header, l1Height uint64) *big.Int {
	// Check if dynamic pricing is enabled at this L1 height
	if !gp.isDynamicPricingEnabledAtHeight(ctx, l1Height) {
		gp.logger.Trace("Using static gas pricing", "minGasPrice", gp.config.BaseFee, "l1Height", l1Height)
		return new(big.Int).Set(gp.config.BaseFee)
	}

	// Dynamic pricing is enabled - use the EIP-1559 calculation
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
	if gp.dynamicPricingEnabled.Load() {
		return gp.config.BaseFee
	}
	// Prior behavior: always use header.BaseFee, even if zero/nil
	if header == nil || header.BaseFee == nil {
		return gp.config.BaseFee
	}
	return header.BaseFee
}
