package crosschain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	busAddress    *common.L1Address
	bridgeAddress *common.L1Address
	storage       storage.Storage
	logger        gethlog.Logger
}

func NewBlockMessageExtractor(
	busAddress *common.L1Address,
	bridgeAddress *common.L1Address,
	storage storage.Storage,
	logger gethlog.Logger,
) BlockMessageExtractor {
	return &blockMessageExtractor{
		busAddress:    busAddress,
		bridgeAddress: bridgeAddress,
		storage:       storage,
		logger:        logger.New(log.CmpKey, log.CrossChainCmp),
	}
}

func (m *blockMessageExtractor) Enabled() bool {
	return m.GetBusAddress().Big().Cmp(gethcommon.Big0) != 0
}

func (m *blockMessageExtractor) StoreCrossChainValueTransfers(ctx context.Context, block *types.Header, processed *common.ProcessedL1Data) error {
	defer core.LogMethodDuration(m.logger, measure.NewStopwatch(), "BlockHeader value transfer messages processed", log.BlockHashKey, block.Hash())

	transferEvents := processed.GetEvents(common.CrossChainMessageTx)
	if len(transferEvents) == 0 {
		return nil
	}

	transfersMessages := make(common.CrossChainMessages, 0)
	var receipts types.Receipts
	for _, txData := range transferEvents {
		if txData.CrossChainMessages == nil {
			continue
		}
		transfersInTx := txData.CrossChainMessages.FilterValueTransfers(*m.bridgeAddress)
		transfersMessages = append(transfersMessages, transfersInTx...)
		receipts = append(receipts, txData.Receipt)
	}

	// Create the ABI components for the ValueTransfer struct
	uint256Type, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)

	// Define the struct components matching the Solidity struct
	valueTransferComponents := abi.Arguments{
		{
			Name: "amount",
			Type: uint256Type,
		},
		{
			Name: "recipient",
			Type: addressType,
		},
	}

	var transfers common.ValueTransferEvents
	for _, msg := range transfersMessages {
		// With publishRawMessage, we can directly unpack the ValueTransfer struct
		// No need to unwrap from CrossChainCall anymore
		unpacked, err := valueTransferComponents.Unpack(msg.Payload)
		if err != nil {
			m.logger.Error("Unable to unpack the value transfer struct", log.ErrKey, err)
			return err
		}

		// Make sure we get the expected number of values
		if len(unpacked) != 2 {
			m.logger.Error("Unexpected number of values unpacked", "expected", 2, "got", len(unpacked))
			return fmt.Errorf("unexpected number of values unpacked: expected 2, got %d", len(unpacked))
		}

		// Convert the unpacked values to the right types
		amount, ok1 := unpacked[0].(*big.Int)
		recipient, ok2 := unpacked[1].(gethcommon.Address)

		if !ok1 || !ok2 {
			m.logger.Error("Failed to convert unpacked values to expected types")
			return fmt.Errorf("failed to convert unpacked values to expected types")
		}

		// Create the ValueTransferEvent with the unpacked data plus message metadata
		valueTransfer := common.ValueTransferEvent{
			Amount:   amount,
			Receiver: recipient,
			Sender:   msg.Sender,
			Sequence: msg.Sequence,
		}

		transfers = append(transfers, valueTransfer)
	}

	err := m.storage.StoreValueTransfers(ctx, block.Hash(), transfers)
	if err != nil {
		m.logger.Error("Unable to store the transfers", log.ErrKey, err)
		return err
	}

	return nil
}

// StoreCrossChainMessages - extracts the cross chain messages for the corresponding block from the receipts.
// The messages will be stored in DB storage for later usage.
// block - the L1 block for which events are extracted.
// processed - all the txs and events relating to the message bus and management contract identified by the logs.
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
			m.logger.Error("Unable to store the messages", log.ErrKey, err)
			return err
		}
	}

	return nil
}

// GetBusAddress - Returns the address of the L1 message bus.
func (m *blockMessageExtractor) GetBusAddress() *common.L1Address {
	return m.busAddress
}
