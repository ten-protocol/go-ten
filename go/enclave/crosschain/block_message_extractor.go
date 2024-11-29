package crosschain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
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
	return m.GetBusAddress().Big().Cmp(gethcommon.Big0) != 0
}

func (m *blockMessageExtractor) StoreCrossChainValueTransfers(ctx context.Context, block *types.Header, processedData *common.ProcessedL1Data) error {
	defer core.LogMethodDuration(m.logger, measure.NewStopwatch(), "BlockHeader value transfer messages processed", log.BlockHashKey, block.Hash())

	// collect all value transfer events from processed data
	var transfers common.ValueTransferEvents
	for _, txData := range processedData.GetEvents(common.CrossChainValueTranserTx) {
		if txData.ValueTransfers != nil {
			transfers = append(transfers, *txData.ValueTransfers...)
		}
	}

	m.logger.Trace("Storing value transfers for block", "nr", len(transfers), log.BlockHashKey, block.Hash())
	err := m.storage.StoreValueTransfers(ctx, block.Hash(), transfers)
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
func (m *blockMessageExtractor) StoreCrossChainMessages(ctx context.Context, block *types.Header, processedData *common.ProcessedL1Data) error {
	defer core.LogMethodDuration(m.logger, measure.NewStopwatch(), "BlockHeader cross chain messages processed", log.BlockHashKey, block.Hash())

	// collect all messages from the events
	var xchain common.CrossChainMessages
	var receipts types.Receipts
	for _, txData := range processedData.GetEvents(common.CrossChainMessageTx) {
		if txData.CrossChainMessages != nil {
			xchain = append(xchain, *txData.CrossChainMessages...)
			receipts = append(receipts, txData.Receipt)
		}
	}

	lazilyLogReceiptChecksum(block, receipts, m.logger)
	if len(xchain) > 0 {
		m.logger.Info(fmt.Sprintf("Storing %d messages for block", len(xchain)), log.BlockHashKey, block.Hash())
		err := m.storage.StoreL1Messages(ctx, block.Hash(), xchain)
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
func (m *blockMessageExtractor) getCrossChainMessages(block *types.Header, receipts common.L1Receipts) (common.CrossChainMessages, error) {
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

	messages, err := ConvertLogsToMessages(logs, CrossChainEventName, MessageBusABI)
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

	transfers, err := ConvertLogsToValueTransfers(logs, ValueTransferEventName, MessageBusABI)
	if err != nil {
		m.logger.Error("Error encountered when converting value transfer receipt logs.", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	return transfers, nil
}
