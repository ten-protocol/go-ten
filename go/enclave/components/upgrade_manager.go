package components

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
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

func (um *upgradeManager) OnL1Block(ctx context.Context, blockHeader *types.Header, processed *common.ProcessedL1Data) error {
	upgradeTxs := processed.GetEvents(common.NetworkUpgradedTx)
	upgrades := make([]NetworkConfig.NetworkConfigUpgraded, 0)
	for _, tx := range upgradeTxs {
		upgrades = append(upgrades, tx.NetworkUpgrades...)
	}
	supported := false
	for _, upgrade := range upgrades {
		um.logger.Info("Upgrade detected", "featureName", upgrade.FeatureName, "featureData", upgrade.FeatureData)
		handlers, ok := um.handlers[upgrade.FeatureName]
		if !ok {
			um.logger.Error("No handler registered for upgrade", "featureName", upgrade.FeatureName)
			return errutil.ErrUpgradeNotSupported
		}
		for _, handler := range handlers {
			canUpgrade := handler.CanUpgrade(ctx, upgrade.FeatureName, upgrade.FeatureData)
			if !canUpgrade {
				um.logger.Error("Upgrade not allowed", "featureName", upgrade.FeatureName)
				return errutil.ErrUpgradeNotSupported
			}
			supported = true
		}
	}
	if !supported {
		return errutil.ErrUpgradeNotSupported
	}
	return nil
}

func (um *upgradeManager) ReplayFinalizedUpgrades(ctx context.Context) error {
	// no-op: upgrades disabled
	return nil
}
