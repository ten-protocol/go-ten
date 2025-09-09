package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"
)

// NetworkUpgradeConfirmationDepth retained for reference
const NetworkUpgradeConfirmationDepth = 64

type upgradeManager struct {
	handlers map[string][]UpgradeHandler
	storage  storage.Storage
	logger   gethlog.Logger
}

func NewUpgradeManager(storage storage.Storage, logger gethlog.Logger) UpgradeManager {
	return &upgradeManager{
		handlers: make(map[string][]UpgradeHandler),
		storage:  storage,
		logger:   logger,
	}
}

func (um *upgradeManager) RegisterUpgradeHandler(featureName string, handler UpgradeHandler) {
	if um.handlers[featureName] == nil {
		um.handlers[featureName] = make([]UpgradeHandler, 0)
	}
	// keep API no-op
}

func (um *upgradeManager) StoreNetworkUpgrades(ctx context.Context, upgrades []NetworkConfig.NetworkConfigUpgraded) error {
	if len(upgrades) == 0 {
		return nil
	}

	for _, upgrade := range upgrades {
		um.logger.Info("Upgrade detected", "featureName", upgrade.FeatureName, "featureData", upgrade.FeatureData)
		err := um.storage.StoreNetworkUpgrade(ctx, &enclavedb.NetworkUpgrade{
			FeatureName: upgrade.FeatureName,
			FeatureData: upgrade.FeatureData,
		})
		if err != nil {
			return fmt.Errorf("failed to store network upgrade. Cause: %w", err)
		}
	}
	um.logger.Info("Stored network upgrades", "upgrades", len(upgrades))
	return nil
}

// collectNetworkUpgrades extracts all NetworkConfigUpgraded from L1 data
func (um *upgradeManager) collectNetworkUpgrades(processed *common.ProcessedL1Data) []NetworkConfig.NetworkConfigUpgraded {
	upgradeTxs := processed.GetEvents(common.NetworkUpgradedTx)
	allUpgrades := make([]NetworkConfig.NetworkConfigUpgraded, 0)
	for _, tx := range upgradeTxs {
		allUpgrades = append(allUpgrades, tx.NetworkUpgrades...)
	}
	return allUpgrades
}

// filterSupportedUpgrades filters upgrades to only include supported ones
func (um *upgradeManager) filterSupportedUpgrades(ctx context.Context, allUpgrades []NetworkConfig.NetworkConfigUpgraded) []NetworkConfig.NetworkConfigUpgraded {
	supportedUpgrades := make([]NetworkConfig.NetworkConfigUpgraded, 0)
	for _, upgrade := range allUpgrades {
		um.logger.Info("Upgrade detected", "featureName", upgrade.FeatureName, "featureData", upgrade.FeatureData)

		// Check if handler is registered for this upgrade
		handlers, ok := um.handlers[upgrade.FeatureName]
		if !ok {
			um.logger.Warn("No handler registered for upgrade, filtering out", "featureName", upgrade.FeatureName)
			continue
		}

		// Check if all handlers support this upgrade
		upgradeSupported := true
		for _, handler := range handlers {
			canUpgrade := handler.CanUpgrade(ctx, upgrade.FeatureName, upgrade.FeatureData)
			if !canUpgrade {
				um.logger.Warn("Upgrade not supported by handler, filtering out", "featureName", upgrade.FeatureName)
				upgradeSupported = false
				break
			}
		}

		if upgradeSupported {
			supportedUpgrades = append(supportedUpgrades, upgrade)
			um.logger.Info("Upgrade supported and included", "featureName", upgrade.FeatureName)
		}
	}

	um.logger.Info("Filtered upgrades", "total", len(allUpgrades), "supported", len(supportedUpgrades))
	return supportedUpgrades
}

func (um *upgradeManager) OnL1Block(ctx context.Context, blockHeader *types.Header, processed *common.ProcessedL1Data) error {
	// Collect all upgrades from L1 data
	allUpgrades := um.collectNetworkUpgrades(processed)

	// Store all upgrades (both supported and unsupported)
	err := um.StoreNetworkUpgrades(ctx, allUpgrades)
	if err != nil {
		um.logger.Error("Failed to store network upgrades", "error", err)
		return err
	}

	// Filter to get only supported upgrades for processing
	supportedUpgrades := um.filterSupportedUpgrades(ctx, allUpgrades)

	// Process only the supported upgrades
	for _, upgrade := range supportedUpgrades {
		um.logger.Info("Processing supported upgrade", "featureName", upgrade.FeatureName, "featureData", upgrade.FeatureData)
	}

	return nil
}

func (um *upgradeManager) ReplayFinalizedUpgrades(ctx context.Context) error {
	// no-op: upgrades disabled
	return nil
}
