package crosschain

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
)

type blockMessageExtractor struct {
	busAddress *common.L1Address
	storage    storage.Storage
	logger     gethlog.Logger
}

func NewBlockMessageExtractor(
	busAddress *common.L1Address,
	storage storage.Storage,
	logger gethlog.Logger,
) BlockMessageExtractor {
	return &blockMessageExtractor{
		busAddress: busAddress,
		storage:    storage,
		logger:     logger.New(log.CmpKey, log.CrossChainCmp),
	}
}

func (m *blockMessageExtractor) Enabled() bool {
	return m.GetBusAddress().Hash().Big().Cmp(gethcommon.Big0) != 0
}

func (m *blockMessageExtractor) StoreCrossChainValueTransfers(block *common.L1Block, receipts common.L1Receipts) error {
	defer m.logger.Info("Block value transfer messages processed", log.BlockHashKey, block.Hash(), log.DurationKey, measure.NewStopwatch())

	/*areReceiptsValid := common.VerifyReceiptHash(block, receipts)

	if !areReceiptsValid && m.Enabled() {
		m.logger.Error("Invalid receipts submitted", log.BlockHashKey, block.Hash())
		return fmt.Errorf("receipts do not match the receipt root for the block")
	}*/

	if len(receipts) == 0 {
		return nil
	}

	transfers, err := m.getValueTransferMessages(receipts)
	if err != nil {
		m.logger.Error("Error encountered while getting inbound value transfers from block", log.BlockHashKey, block.Hash(), log.ErrKey, err)
		return err
	}

	hasTransfers := len(transfers) > 0
	if !hasTransfers {
		return nil
	}

	m.logger.Trace(fmt.Sprintf("Storing %d value transfers for block", len(transfers)), log.BlockHashKey, block.Hash())
	err = m.storage.StoreValueTransfers(block.Hash(), transfers)
	if err != nil {
		m.logger.Crit("Unable to store the transfers", log.ErrKey, err)
		return err
	}

	return nil
}

// StoreCrossChainMessages - extracts the cross chain messages for the corresponding block from the receipts.
// The messages will be stored in DB storage for later usage.
// block - the L1 block for which events are extracted.
// receipts - all of the receipts for the corresponding block. This is validated.
func (m *blockMessageExtractor) StoreCrossChainMessages(block *common.L1Block, receipts common.L1Receipts) error {
	defer m.logger.Info("Block cross chain messages processed", log.BlockHashKey, block.Hash(), log.DurationKey, measure.NewStopwatch())

	if len(receipts) == 0 {
		return nil
	}

	lazilyLogReceiptChecksum(fmt.Sprintf("Processing block: %s receipts: %d", block.Hash(), len(receipts)), receipts, m.logger)
	messages, err := m.getCrossChainMessages(block, receipts)
	if err != nil {
		m.logger.Error("Converting receipts to messages failed.", log.ErrKey, err)
		return err
	}

	if len(messages) > 0 {
		m.logger.Info(fmt.Sprintf("Storing %d messages for block", len(messages)), log.BlockHashKey, block.Hash())
		err = m.storage.StoreL1Messages(block.Hash(), messages)
		if err != nil {
			m.logger.Crit("Unable to store the messages", log.ErrKey, err)
			return err
		}
	}

	return nil
}

// GetBusAddress - Returns the address of the L1 message bus.
func (m *blockMessageExtractor) GetBusAddress() *common.L1Address {
	return m.busAddress
}

// getCrossChainMessages - Converts the relevant logs from the appropriate message bus address to synthetic transactions and returns them
func (m *blockMessageExtractor) getCrossChainMessages(block *common.L1Block, receipts common.L1Receipts) (common.CrossChainMessages, error) {
	if len(receipts) == 0 {
		return make(common.CrossChainMessages, 0), nil
	}

	// Retrieves the relevant logs from the message bus.
	logs, err := filterLogsFromReceipts(receipts, m.GetBusAddress(), &CrossChainEventID)
	if err != nil {
		m.logger.Error("Error encountered when filtering receipt logs.", log.ErrKey, err)
		return make(common.CrossChainMessages, 0), err
	}
	m.logger.Trace("Extracted cross chain logs from receipts", "logCount", len(logs))

	messages, err := convertLogsToMessages(logs, CrossChainEventName, MessageBusABI)
	if err != nil {
		m.logger.Error("Error encountered converting the extracted relevant logs to messages", log.ErrKey, err)
		return make(common.CrossChainMessages, 0), err
	}

	m.logger.Trace(fmt.Sprintf("Found %d cross chain messages that will be submitted to L2!", len(messages)), log.BlockHashKey, block.Hash())

	return messages, nil
}

func (m *blockMessageExtractor) getValueTransferMessages(receipts common.L1Receipts) (common.ValueTransferEvents, error) {
	if len(receipts) == 0 {
		return make(common.ValueTransferEvents, 0), nil
	}

	// Retrieves the relevant logs from the message bus.
	logs, err := filterLogsFromReceipts(receipts, m.GetBusAddress(), &ValueTransferEventID)
	if err != nil {
		m.logger.Error("Error encountered when filtering receipt logs.", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	transfers, err := convertLogsToValueTransfers(logs, ValueTransferEventName, MessageBusABI)
	if err != nil {
		m.logger.Error("Error encountered when converting value transfer receipt logs.", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	return transfers, nil
}
