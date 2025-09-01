package components

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
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
	// no-op: upgrades disabled
	return nil
}

func (um *upgradeManager) ReplayFinalizedUpgrades(ctx context.Context) error {
	// no-op: upgrades disabled
	return nil
}
