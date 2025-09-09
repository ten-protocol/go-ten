package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
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

func (um *upgradeManager) StoreNetworkUpgrades(ctx context.Context, blockHeader *types.Header, upgrades []NetworkConfig.NetworkConfigUpgraded) error {
	if len(upgrades) == 0 {
		return nil
	}

	for _, upgrade := range upgrades {
		um.logger.Info("Upgrade detected", "featureName", upgrade.FeatureName, "featureData", upgrade.FeatureData)
		blockHeight := blockHeader.Number.Uint64() + 64
		err := um.storage.StoreNetworkUpgrade(ctx, &enclavedb.NetworkUpgrade{
			FeatureName:       upgrade.FeatureName,
			FeatureData:       upgrade.FeatureData,
			BlockHash:         blockHeader.Hash(),
			BlockHeightFinal:  &blockHeight,
			BlockHeightActive: &blockHeight,
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

func (um *upgradeManager) getUpgradesFor(blockHeader *types.Header) ([]NetworkConfig.NetworkConfigUpgraded, error) {
	ctx := context.Background()

	// Ensure the current block is canonical before proceeding
	isCanonical, err := um.storage.IsBlockCanonical(ctx, blockHeader.Hash())
	if err != nil {
		um.logger.Error("Failed to check block canonicality", "hash", blockHeader.Hash(), "error", err)
		return nil, fmt.Errorf("failed to check block canonicality: %w", err)
	}
	if !isCanonical {
		return nil, fmt.Errorf("block %s is not canonical", blockHeader.Hash())
	}

	// Load activated upgrades up to current block height
	currentHeight := blockHeader.Number.Uint64()
	stored, err := um.storage.GetActivatedNetworkUpgrades(ctx, currentHeight)
	if err != nil {
		um.logger.Error("Failed to load activated network upgrades from storage", "error", err, "height", currentHeight)
		return nil, fmt.Errorf("failed to load activated network upgrades: %w", err)
	}

	// Filter: only upgrades whose saved block hash is canonical
	candidates := make([]NetworkConfig.NetworkConfigUpgraded, 0)
	for _, u := range stored {
		// The upgrade must have been recorded for a canonical block hash
		upgradeBlockCanonical, err := um.storage.IsBlockCanonical(ctx, u.BlockHash)
		if err != nil {
			um.logger.Error("Failed to check canonicality for upgrade's block", "hash", u.BlockHash, "error", err)
			return nil, fmt.Errorf("failed to check canonicality for upgrade's block: %w", err)
		}
		if !upgradeBlockCanonical {
			continue
		}

		candidates = append(candidates, NetworkConfig.NetworkConfigUpgraded{
			FeatureName: u.FeatureName,
			FeatureData: u.FeatureData,
		})
	}

	// Verify all selected upgrades are supported by registered handlers
	supported := um.filterSupportedUpgrades(ctx, candidates)
	if len(supported) != len(candidates) {
		return nil, fmt.Errorf("unsupported upgrades detected at height %d: total=%d supported=%d", currentHeight, len(candidates), len(supported))
	}

	return supported, nil
}

func (um *upgradeManager) OnL1Block(ctx context.Context, blockHeader *types.Header, processed *common.ProcessedL1Data) error {
	// Collect all upgrades from L1 data
	allUpgrades := um.collectNetworkUpgrades(processed)

	// Store all upgrades (both supported and unsupported)
	err := um.StoreNetworkUpgrades(ctx, blockHeader, allUpgrades)
	if err != nil {
		um.logger.Error("Failed to store network upgrades", "error", err)
		return err
	}

	_, err = um.getUpgradesFor(blockHeader)
	if err != nil {
		um.logger.Error("Failed to get upgrades for block", "error", err)
		return errutil.ErrUpgradeNotSupported
	}

	return nil
}

func (um *upgradeManager) ReplayFinalizedUpgrades(ctx context.Context) error {
	// no-op: upgrades disabled
	return nil
}
