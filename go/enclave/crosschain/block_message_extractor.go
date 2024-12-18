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

func (m *blockMessageExtractor) StoreCrossChainValueTransfers(ctx context.Context, block *types.Header, processed *common.ProcessedL1Data) error {
	defer core.LogMethodDuration(m.logger, measure.NewStopwatch(), "BlockHeader value transfer messages processed", log.BlockHashKey, block.Hash())

	transferEvents := processed.GetEvents(common.CrossChainValueTranserTx)
	if len(transferEvents) == 0 {
		return nil
	}

	var transfers common.ValueTransferEvents
	for _, txData := range transferEvents {
		if txData.ValueTransfers != nil {
			transfers = append(transfers, txData.ValueTransfers...)
		}
	}

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
// processed - all the txs and events relating to the messagebus identified by the logs. This is validated.
func (m *blockMessageExtractor) StoreCrossChainMessages(ctx context.Context, block *types.Header, processed *common.ProcessedL1Data) error {
	defer core.LogMethodDuration(m.logger, measure.NewStopwatch(), "BlockHeader cross chain messages processed", log.BlockHashKey, block.Hash())

	messageEvents := processed.GetEvents(common.CrossChainMessageTx)
	if len(messageEvents) == 0 {
		return nil
	}
	var messages common.CrossChainMessages
	var receipts types.Receipts
	for _, txData := range messageEvents {
		if txData.CrossChainMessages != nil {
			messages = append(messages, txData.CrossChainMessages...)
			receipts = append(receipts, txData.Receipt)
		}
	}

	lazilyLogReceiptChecksum(block, receipts, m.logger)
	if len(messages) > 0 {
		m.logger.Info(fmt.Sprintf("Storing %d messages for block", len(messages)), log.BlockHashKey, block.Hash())
		err := m.storage.StoreL1Messages(ctx, block.Hash(), messages)
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
