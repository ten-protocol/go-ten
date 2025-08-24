package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

// NetworkUpgradeConfirmationDepth is the number of L1 blocks required for finality
// 64 blocks = 2 epochs ≈ 12.8 minutes on Ethereum PoS
const NetworkUpgradeConfirmationDepth = 64

type upgradeManager struct {
	// Map of featureName -> list of handlers
	handlers map[string][]UpgradeHandler

	// Database storage for network upgrades
	storage storage.Storage

	logger gethlog.Logger
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
	um.handlers[featureName] = append(um.handlers[featureName], handler)
	um.logger.Info("Registered upgrade handler", "featureName", featureName)
}

func (um *upgradeManager) OnL1Block(ctx context.Context, blockHeader *types.Header, processed *common.ProcessedL1Data) error {
	// First, add any new network upgrade events as pending
	err := um.addNewPendingUpgrades(ctx, blockHeader, processed)
	if err != nil {
		return fmt.Errorf("failed to add new pending upgrades: %w", err)
	}

	// Then, check if any pending upgrades have reached finality
	err = um.checkFinalityAndApplyUpgrades(ctx, blockHeader)
	if err != nil {
		return fmt.Errorf("failed to check finality and apply upgrades: %w", err)
	}

	return nil
}

func (um *upgradeManager) addNewPendingUpgrades(ctx context.Context, blockHeader *types.Header, processed *common.ProcessedL1Data) error {
	networkUpgradeTxs := processed.GetEvents(common.NetworkUpgradedTx)
	if len(networkUpgradeTxs) == 0 {
		return nil
	}

	for _, txData := range networkUpgradeTxs {
		if len(txData.NetworkUpgrades) == 0 {
			continue
		}

		for _, upgrade := range txData.NetworkUpgrades {
			// Check if any registered handlers can process this upgrade
			handlers := um.handlers[upgrade.FeatureName]
			if len(handlers) == 0 {
				um.logger.Error("No handlers registered for network upgrade feature - node restart required",
					"featureName", upgrade.FeatureName,
					"txHash", txData.Transaction.Hash().Hex())
				return fmt.Errorf("no handlers registered for network upgrade feature '%s' - restart node with newer version", upgrade.FeatureName)
			}

			// Check if at least one handler can process this upgrade
			canProcess := false
			for _, handler := range handlers {
				if handler.CanUpgrade(ctx, upgrade.FeatureName, upgrade.FeatureData) {
					canProcess = true
					break
				}
			}

			if !canProcess {
				um.logger.Error("Network upgrade not supported by current node version - restart required",
					"featureName", upgrade.FeatureName,
					"featureData", string(upgrade.FeatureData),
					"txHash", txData.Transaction.Hash().Hex())
				return fmt.Errorf("network upgrade '%s' with data '%s' not supported by current node version - restart node with newer version",
					upgrade.FeatureName, string(upgrade.FeatureData))
			}

			// All checks passed - store pending upgrade in database
			err := um.storage.StorePendingNetworkUpgrade(
				ctx,
				upgrade.FeatureName,
				upgrade.FeatureData,
				blockHeader.Number.Uint64(),
				blockHeader.Hash(),
				txData.Transaction.Hash(),
			)
			if err != nil {
				return fmt.Errorf("failed to store pending network upgrade: %w", err)
			}

			um.logger.Info("Added pending network upgrade",
				"featureName", upgrade.FeatureName,
				"dataLength", len(upgrade.FeatureData),
				"l1Height", blockHeader.Number.Uint64(),
				"l1Hash", blockHeader.Hash().Hex(),
				"txHash", txData.Transaction.Hash().Hex())
		}
	}

	return nil
}

func (um *upgradeManager) checkFinalityAndApplyUpgrades(ctx context.Context, blockHeader *types.Header) error {
	currentHeight := blockHeader.Number.Uint64()

	// Get all pending upgrades from database
	pendingUpgrades, err := um.storage.GetPendingNetworkUpgrades(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending upgrades: %w", err)
	}

	// Check each pending upgrade for finality
	for _, dbUpgrade := range pendingUpgrades {
		// Check if this upgrade has reached finality
		if currentHeight >= dbUpgrade.AppliedAtL1Height+NetworkUpgradeConfirmationDepth {
			// Upgrade has reached finality - finalize it in database
			err := um.storage.FinalizeNetworkUpgrade(ctx, dbUpgrade.TxHash, currentHeight, blockHeader.Hash())
			if err != nil {
				return fmt.Errorf("failed to finalize network upgrade: %w", err)
			}

			um.logger.Info("Network upgrade reached finality",
				"featureName", dbUpgrade.FeatureName,
				"originalHeight", dbUpgrade.AppliedAtL1Height,
				"finalizedAtHeight", currentHeight,
				"confirmations", currentHeight-dbUpgrade.AppliedAtL1Height)

			// Apply to all registered handlers
			err = um.notifyHandlers(ctx, dbUpgrade.FeatureName, dbUpgrade.FeatureData)
			if err != nil {
				um.logger.Error("Failed to notify handlers for finalized upgrade",
					"featureName", dbUpgrade.FeatureName,
					"error", err)
				return fmt.Errorf("failed to notify handlers for feature %s: %w", dbUpgrade.FeatureName, err)
			}

			um.logger.Info("Successfully applied finalized network upgrade",
				"featureName", dbUpgrade.FeatureName)
		}
	}

	return nil
}

func (um *upgradeManager) notifyHandlers(ctx context.Context, featureName string, featureData []byte) error {
	handlers := um.handlers[featureName]
	if len(handlers) == 0 {
		um.logger.Warn("No handlers registered for feature", "featureName", featureName)
		return nil
	}

	for _, handler := range handlers {
		err := handler.HandleUpgrade(ctx, featureName, featureData)
		if err != nil {
			um.logger.Error("Handler failed to process upgrade",
				"featureName", featureName,
				"error", err,
				log.ErrKey, err)
			return fmt.Errorf("handler failed to process upgrade for feature %s: %w", featureName, err)
		}
	}

	um.logger.Info("Successfully notified all handlers",
		"featureName", featureName,
		"handlersNotified", len(handlers))

	return nil
}

func (um *upgradeManager) ReplayFinalizedUpgrades(ctx context.Context) error {
	// Get all finalized upgrades from database
	finalizedUpgrades, err := um.storage.GetFinalizedNetworkUpgrades(ctx)
	if err != nil {
		return fmt.Errorf("failed to get finalized upgrades: %w", err)
	}

	um.logger.Info("Replaying finalized upgrades on startup", "count", len(finalizedUpgrades))

	for _, dbUpgrade := range finalizedUpgrades {
		err := um.notifyHandlers(ctx, dbUpgrade.FeatureName, dbUpgrade.FeatureData)
		if err != nil {
			um.logger.Error("Failed to replay finalized upgrade",
				"featureName", dbUpgrade.FeatureName,
				"error", err)
			return fmt.Errorf("failed to replay finalized upgrade for feature %s: %w", dbUpgrade.FeatureName, err)
		}
	}

	um.logger.Info("Successfully replayed all finalized upgrades")
	return nil
}
