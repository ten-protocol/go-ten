package crosschain

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/holiman/uint256"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

type MessageBusManager struct {
	messageBusAddress *gethcommon.Address
	bridgeAddress     *gethcommon.Address
	storage           storage.Storage
	logger            gethlog.Logger
}

func NewTenMessageBusManager(
	storage storage.Storage,
	logger gethlog.Logger,
) Manager {
	logger = logger.New(log.CmpKey, log.CrossChainCmp)

	return &MessageBusManager{
		messageBusAddress: nil,
		bridgeAddress:     nil,
		storage:           storage,
		logger:            logger,
	}
}

// GetBusAddress - Returns the L2 address of the message bus contract.
// todo - figure out how to expose the deployed contract to the external world. Perhaps extract event from contract construction?
func (m *MessageBusManager) GetBusAddress() *common.L2Address {
	return m.messageBusAddress
}

// Initialize - Derives the address of the message bus contract.
func (m *MessageBusManager) Initialize(systemAddresses common.SystemContractAddresses) error {
	address, ok := systemAddresses["MessageBus"]
	if !ok {
		return fmt.Errorf("message bus contract not found in system addresses")
	}

	bridgeAddress, ok := systemAddresses["EthereumBridge"]
	if !ok {
		return fmt.Errorf("ethereum bridge contract not found in system addresses")
	}

	m.messageBusAddress = address
	m.bridgeAddress = bridgeAddress
	return nil
}

// ExtractOutboundMessages - Finds relevant logs in the receipts and converts them to cross chain messages.
func (m *MessageBusManager) ExtractOutboundMessages(_ context.Context, receipts common.L2Receipts) (common.CrossChainMessages, error) {
	logs, err := filterLogsFromReceipts(receipts, m.messageBusAddress, &ethadapter.CrossChainEventID)
	if err != nil {
		m.logger.Error("Error extracting logs from L2 message bus!", log.ErrKey, err)
		return make(common.CrossChainMessages, 0), err
	}

	messages, err := ConvertLogsToMessages(logs, ethadapter.CrossChainEventName, ethadapter.MessageBusABI)
	if err != nil {
		m.logger.Error("Error converting messages from L2 message bus!", log.ErrKey, err)
		return make(common.CrossChainMessages, 0), err
	}

	return messages, nil
}

// ExtractOutboundTransfers - Finds relevant logs in the receipts and converts them to cross chain messages.
func (m *MessageBusManager) ExtractOutboundTransfers(_ context.Context, receipts common.L2Receipts) (common.ValueTransferEvents, error) {
	logs, err := filterLogsFromReceipts(receipts, m.messageBusAddress, &ethadapter.CrossChainEventID)
	if err != nil {
		m.logger.Error("Error extracting logs from L2 message bus!", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	msgs, err := ConvertLogsToMessages(logs, ethadapter.CrossChainEventName, ethadapter.MessageBusABI)
	if err != nil {
		m.logger.Error("Error converting messages from L2 message bus!", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	transfers, err := ConvertCrossChainMessagesToValueTransfers(msgs, ethadapter.ValueTransferEventName, m.bridgeAddress)
	if err != nil {
		m.logger.Error("Error converting transfers from L2 message bus!", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	return transfers, nil
}

// RetrieveInboundMessages - Retrieves the cross chain messages between two blocks.
// todo (@stefan) - fix ordering of messages, currently it is irrelevant.
// todo (@stefan) - do not extract messages below their consistency level. Irrelevant security wise.
// todo (@stefan) - surface errors
func (m *MessageBusManager) RetrieveInboundMessages(ctx context.Context, fromBlock *types.Header, toBlock *types.Header) (common.CrossChainMessages, common.ValueTransferEvents, error) {
	messages := make(common.CrossChainMessages, 0)
	transfers := make(common.ValueTransferEvents, 0)

	from := fromBlock.Hash()
	height := fromBlock.Number.Uint64()
	if !m.storage.IsAncestor(ctx, toBlock, fromBlock) {
		m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
	}
	// Iterate through the blocks.
	b := toBlock
	for {
		if bytes.Equal(b.Hash().Bytes(), from.Bytes()) {
			break
		}

		m.logger.Trace(fmt.Sprintf("Looking for cross chain messages at block %s", b.Hash().Hex()))

		messagesForBlock, err := m.storage.GetL1Messages(ctx, b.Hash())
		if err != nil {
			return nil, nil, fmt.Errorf("reading the key for the block failed with uncommon reason: %w", err)
		}

		transfersForBlock, err := m.storage.GetL1Transfers(ctx, b.Hash())
		if err != nil {
			return nil, nil, fmt.Errorf("unable to get L1 transfers for block that should be there %w", err)
		}

		messages = append(messages, messagesForBlock...) // Ordering here might work in POBI, but might be weird for fast finality
		transfers = append(transfers, transfersForBlock...)

		// No deposits before genesis.
		if b.Number.Uint64() < height {
			return nil, nil, fmt.Errorf("block height is less than genesis height")
		}
		p, err := m.storage.FetchBlock(ctx, b.ParentHash)
		if err != nil {
			return nil, nil, fmt.Errorf("synthetic transactions can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}

	logf := m.logger.Info
	if len(messages)+len(transfers) == 0 {
		logf = m.logger.Debug
	}
	logf(fmt.Sprintf("Extracted cross chain messages for block height %d ->%d", fromBlock.Number.Uint64(), toBlock.Number.Uint64()), "no_msgs", len(messages), "no_value_transfers", len(transfers))

	return messages, transfers, nil
}

const BalanceIncreaseXChainValueTransfer tracing.BalanceChangeReason = 110

func (m *MessageBusManager) ExecuteValueTransfers(_ context.Context, transfers common.ValueTransferEvents, rollupState *state.StateDB) {
	for _, transfer := range transfers {
		rollupState.AddBalance(transfer.Receiver, uint256.MustFromBig(transfer.Amount), BalanceIncreaseXChainValueTransfer)
		m.logger.Debug(fmt.Sprintf("Executed cross chain value transfer from %s to %s with amount %s", transfer.Sender.Hex(), transfer.Receiver.Hex(), transfer.Amount.String()))
	}
}

// CreateSyntheticTransactions - generates transactions that the enclave should execute internally for the messages.
func (m *MessageBusManager) CreateSyntheticTransactions(_ context.Context, messages common.CrossChainMessages, transfers common.ValueTransferEvents, rollupState *state.StateDB) (common.L2Transactions, error) {
	if len(messages) == 0 && len(transfers) == 0 {
		return make(common.L2Transactions, 0), nil
	}

	if m.messageBusAddress == nil {
		m.logger.Crit("Message bus address not set")
	}

	// Get current nonce for this stateDB.
	// There can be forks thus we cannot trust the wallet.
	startingNonce := rollupState.GetNonce(common.MaskedSender(*m.messageBusAddress))

	syntheticTransactions := make(types.Transactions, 0)
	for idx, message := range messages {
		delayInBlocks := big.NewInt(int64(message.ConsistencyLevel))
		data, err := ethadapter.MessageBusABI.Pack("storeCrossChainMessage", message, delayInBlocks)
		if err != nil {
			return nil, fmt.Errorf("failed packing storeCrossChainMessage %w", err)
		}

		tx := &types.LegacyTx{
			Nonce:    startingNonce + uint64(idx),
			Value:    gethcommon.Big0,
			Gas:      5_000_000,
			GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
			Data:     data,
			To:       m.messageBusAddress,
		}

		stx := types.NewTx(tx)
		syntheticTransactions = append(syntheticTransactions, stx)
	}

	startingNonce += uint64(len(messages))

	for idx, transfer := range transfers {
		data, err := ethadapter.MessageBusABI.Pack("notifyDeposit", transfer.Receiver, transfer.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed packing notifyDeposit %w", err)
		}

		tx := &types.LegacyTx{
			Nonce:    startingNonce + uint64(idx),
			Value:    gethcommon.Big0,
			Data:     data,
			To:       m.messageBusAddress,
			Gas:      5_000_000,
			GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		}
		syntheticTransactions = append(syntheticTransactions, types.NewTx(tx))
	}

	return syntheticTransactions, nil
}

// ConvertLogsToValueTransfers converts logs from the message bus into value transfer events
func ConvertCrossChainMessagesToValueTransfers(msgs common.CrossChainMessages, eventName string, bridgeAddress *common.L1Address) (common.ValueTransferEvents, error) {

	transfers := msgs.FilterValueTransfers(*bridgeAddress)

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

	valueTransfers := make(common.ValueTransferEvents, 0)
	for _, transfer := range transfers {
		// Unpack the log data
		unpacked, err := valueTransferComponents.Unpack(transfer.Payload)
		if err != nil {
			return nil, fmt.Errorf("unable to unpack the value transfer struct: %w", err)
		}

		// Make sure we get the expected number of values
		if len(unpacked) != 2 {
			return nil, fmt.Errorf("unexpected number of values unpacked: expected 2, got %d", len(unpacked))
		}

		// Convert the unpacked values to the right types
		amount, ok1 := unpacked[0].(*big.Int)
		recipient, ok2 := unpacked[1].(gethcommon.Address)

		if !ok1 || !ok2 {
			return nil, fmt.Errorf("failed to convert unpacked values to expected types")
		}

		// Create the ValueTransferEvent with the unpacked data
		valueTransfer := common.ValueTransferEvent{
			Amount:   amount,
			Receiver: recipient,
			Sender:   gethcommon.BytesToAddress(transfer.Sender.Bytes()), // Sender is in the first topic
		}

		valueTransfers = append(valueTransfers, valueTransfer)
	}

	return valueTransfers, nil
}
