package crosschain

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
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

func (m *MessageBusManager) GetL2BridgeAddress() *common.L2Address {
	return m.bridgeAddress
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
			Gas:      params.MaxTxGas,
			GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
			Data:     data,
			To:       m.messageBusAddress,
		}

		stx := types.NewTx(tx)
		syntheticTransactions = append(syntheticTransactions, stx)
	}

	return syntheticTransactions, nil
}

// ConvertLogsToValueTransfers converts logs from the message bus into value transfer events
func ConvertCrossChainMessagesToValueTransfers(msgs common.CrossChainMessages, eventName string, bridgeAddress *common.L1Address) (common.ValueTransferEvents, error) {
	transfers := msgs.FilterValueTransfers(*bridgeAddress)

	valueTransfers := make(common.ValueTransferEvents, 0)
	for _, m := range transfers {
		// 1) Decode CrossChainCall wrapper (tuple)
		out, err := decodeCrossChainCall(m.Payload)
		if err != nil {
			return nil, fmt.Errorf("unpack CrossChainCall failed: topic=%d sender=%s seq=%d: %w", m.Topic, m.Sender.Hex(), m.Sequence, err)
		}

		// 2) Decode inner receiveAssets call
		asset, amount, recipient, err := decodeReceiveAssetsCall(out.Data)
		if err != nil {
			selector := ""
			if len(out.Data) >= 4 {
				selector = gethcommon.Bytes2Hex(out.Data[:4])
			}
			return nil, fmt.Errorf("unable to unpack receiveAssets call data: topic=%d sender=%s seq=%d selector=%s: %w", m.Topic, m.Sender.Hex(), m.Sequence, selector, err)
		}

		// Only native asset transfers (asset == address(0)) affect L2 native balance
		if asset == (gethcommon.Address{}) {
			valueTransfers = append(valueTransfers, common.ValueTransferEvent{
				Amount:   amount,
				Receiver: recipient,
				Sender:   m.Sender,
			})
		}
	}

	return valueTransfers, nil
}

// crossChainCall mirrors the Solidity ICrossChainMessenger.CrossChainCall struct
type crossChainCall struct {
	Target gethcommon.Address
	Data   []byte
	Gas    *big.Int
}

func decodeCrossChainCall(payload []byte) (*crossChainCall, error) {
	tupleType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "target", Type: "address"},
		{Name: "data", Type: "bytes"},
		{Name: "gas", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("build tuple type: %w", err)
	}
	args := abi.Arguments{{Type: tupleType}}
	unpacked, err := args.Unpack(payload)
	if err != nil {
		head := 32
		if len(payload) < head {
			head = len(payload)
		}
		return nil, fmt.Errorf("len=%d head=%s: %w", len(payload), gethcommon.Bytes2Hex(payload[:head]), err)
	}
	if len(unpacked) != 1 {
		return nil, fmt.Errorf("unexpected tuple arity: got=%d", len(unpacked))
	}
	conv := abi.ConvertType(unpacked[0], new(crossChainCall))
	out, ok := conv.(*crossChainCall)
	if !ok {
		return nil, fmt.Errorf("failed to convert tuple to crossChainCall")
	}
	if out.Target == (gethcommon.Address{}) {
		return nil, fmt.Errorf("cross chain call has zero target")
	}
	return out, nil
}

func decodeReceiveAssetsCall(data []byte) (gethcommon.Address, *big.Int, gethcommon.Address, error) {
	if len(data) < 4 {
		return gethcommon.Address{}, nil, gethcommon.Address{}, fmt.Errorf("calldata too short for selector")
	}
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	args := abi.Arguments{{Type: addressType}, {Type: uint256Type}, {Type: addressType}}
	params, err := args.Unpack(data[4:])
	if err != nil {
		return gethcommon.Address{}, nil, gethcommon.Address{}, err
	}
	if len(params) != 3 {
		return gethcommon.Address{}, nil, gethcommon.Address{}, fmt.Errorf("unexpected arg count: %d", len(params))
	}
	asset, ok1 := params[0].(gethcommon.Address)
	amount, ok2 := params[1].(*big.Int)
	receiver, ok3 := params[2].(gethcommon.Address)
	if !ok1 || !ok2 || !ok3 {
		return gethcommon.Address{}, nil, gethcommon.Address{}, fmt.Errorf("type assertion failed")
	}
	return asset, amount, receiver, nil
}
