package crosschain

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/holiman/uint256"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/wallet"
)

// todo (#1549) - implement with cryptography epic
const ( // DO NOT USE OR CHANGE THIS KEY IN THE REST OF THE CODEBASE
	ownerKeyHex = "6e384a07a01263518a18a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
)

type MessageBusManager struct {
	messageBusAddress *gethcommon.Address
	storage           storage.Storage
	logger            gethlog.Logger
	wallet            wallet.Wallet
}

func NewObscuroMessageBusManager(
	storage storage.Storage, /*key *ecdsa.PrivateKey,*/
	chainID *big.Int,
	logger gethlog.Logger,
) Manager {
	// todo (#1549) - implement with cryptography epic, remove this key and use the DeriveKey
	key, _ := crypto.HexToECDSA(ownerKeyHex)
	logger = logger.New(log.CmpKey, log.CrossChainCmp)
	wallet := wallet.NewInMemoryWalletFromPK(chainID, key, logger)

	logger.Info(fmt.Sprintf("L2 Cross Chain Owner Address: %s", wallet.Address().Hex()))

	return &MessageBusManager{
		messageBusAddress: nil,
		storage:           storage,
		logger:            logger,
		wallet:            wallet,
	}
}

func (m *MessageBusManager) IsSyntheticTransaction(transaction *common.L2Tx) bool {
	sender, err := core.GetExternalTxSigner(transaction)
	if err != nil {
		return false
	}
	// The message bus manager considers the transaction synthetic only if the sender is
	// the owner identity which should be available only to enclaves.
	return bytes.Equal(sender.Bytes(), m.GetOwner().Bytes())
}

// GetOwner - Returns the address of the identity owning the message bus.
func (m *MessageBusManager) GetOwner() common.L2Address {
	return m.wallet.Address()
}

// GetBusAddress - Returns the L2 address of the message bus contract.
// todo - figure out how to expose the deployed contract to the external world. Perhaps extract event from contract construction?
func (m *MessageBusManager) GetBusAddress() *common.L2Address {
	return m.messageBusAddress
}

// DeriveMessageBusAddress - Derives the address of the message bus contract.
func (m *MessageBusManager) Initialize(systemAddresses common.SystemContractAddresses) error {
	address, ok := systemAddresses["MessageBus"]
	if !ok {
		return fmt.Errorf("message bus contract not found in system addresses")
	}

	m.messageBusAddress = address
	return nil
}

// GenerateMessageBusDeployTx - Returns a signed message bus deployment transaction.
func (m *MessageBusManager) GenerateMessageBusDeployTx() (*common.L2Tx, error) {
	tx := &types.LegacyTx{
		Nonce:    0, // The first transaction of the owner identity should always be deploying the contract
		Value:    gethcommon.Big0,
		Gas:      5_000_000,       // It's quite the expensive contract.
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     gethcommon.FromHex(MessageBus.MessageBusMetaData.Bin),
		To:       nil, // Geth requires nil instead of gethcommon.Address{} which equates to zero address in order to return receipt.
	}

	stx, err := m.wallet.SignTransaction(tx)
	if err != nil {
		return nil, err
	}

	m.logger.Trace(fmt.Sprintf("Generated synthetic deployment transaction for the MessageBus contract %s - TX HASH: %s", m.messageBusAddress.Hex(), stx.Hash().Hex()),
		log.CmpKey, log.CrossChainCmp)

	return stx, nil
}

// ExtractOutboundMessages - Finds relevant logs in the receipts and converts them to cross chain messages.
func (m *MessageBusManager) ExtractOutboundMessages(ctx context.Context, receipts common.L2Receipts) (common.CrossChainMessages, error) {
	logs, err := filterLogsFromReceipts(receipts, m.messageBusAddress, &CrossChainEventID)
	if err != nil {
		m.logger.Error("Error extracting logs from L2 message bus!", log.ErrKey, err)
		return make(common.CrossChainMessages, 0), err
	}

	messages, err := ConvertLogsToMessages(logs, CrossChainEventName, MessageBusABI)
	if err != nil {
		m.logger.Error("Error converting messages from L2 message bus!", log.ErrKey, err)
		return make(common.CrossChainMessages, 0), err
	}

	return messages, nil
}

// ExtractOutboundTransfers - Finds relevant logs in the receipts and converts them to cross chain messages.
func (m *MessageBusManager) ExtractOutboundTransfers(_ context.Context, receipts common.L2Receipts) (common.ValueTransferEvents, error) {
	logs, err := filterLogsFromReceipts(receipts, m.messageBusAddress, &ValueTransferEventID)
	if err != nil {
		m.logger.Error("Error extracting logs from L2 message bus!", log.ErrKey, err)
		return make(common.ValueTransferEvents, 0), err
	}

	transfers, err := ConvertLogsToValueTransfers(logs, ValueTransferEventName, MessageBusABI)
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
func (m *MessageBusManager) RetrieveInboundMessages(ctx context.Context, fromBlock *types.Header, toBlock *types.Header) (common.CrossChainMessages, common.ValueTransferEvents) {
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
			m.logger.Crit("Reading the key for the block failed with uncommon reason.", log.ErrKey, err)
		}

		transfersForBlock, err := m.storage.GetL1Transfers(ctx, b.Hash())
		if err != nil {
			m.logger.Crit("Unable to get L1 transfers for block that should be there.", log.ErrKey, err)
		}

		messages = append(messages, messagesForBlock...) // Ordering here might work in POBI, but might be weird for fast finality
		transfers = append(transfers, transfersForBlock...)

		// No deposits before genesis.
		if b.Number.Uint64() < height {
			m.logger.Crit("block height is less than genesis height")
		}
		p, err := m.storage.FetchBlock(ctx, b.ParentHash)
		if err != nil {
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}

	logf := m.logger.Info
	if len(messages)+len(transfers) == 0 {
		logf = m.logger.Debug
	}
	logf(fmt.Sprintf("Extracted cross chain messages for block height %d ->%d", fromBlock.Number.Uint64(), toBlock.Number.Uint64()), "no_msgs", len(messages), "no_value_transfers", len(transfers))

	return messages, transfers
}

const BalanceIncreaseXChainValueTransfer tracing.BalanceChangeReason = 110

func (m *MessageBusManager) ExecuteValueTransfers(ctx context.Context, transfers common.ValueTransferEvents, rollupState *state.StateDB) {
	for _, transfer := range transfers {
		rollupState.AddBalance(transfer.Receiver, uint256.MustFromBig(transfer.Amount), BalanceIncreaseXChainValueTransfer)
		m.logger.Debug(fmt.Sprintf("Executed cross chain value transfer from %s to %s with amount %s", transfer.Sender.Hex(), transfer.Receiver.Hex(), transfer.Amount.String()))
	}
}

// CreateSyntheticTransactions - generates transactions that the enclave should execute internally for the messages.
func (m *MessageBusManager) CreateSyntheticTransactions(ctx context.Context, messages common.CrossChainMessages, transfers common.ValueTransferEvents, rollupState *state.StateDB) common.L2Transactions {
	if len(messages) == 0 && len(transfers) == 0 {
		return make(common.L2Transactions, 0)
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
		data, err := MessageBusABI.Pack("storeCrossChainMessage", message, delayInBlocks)
		if err != nil {
			m.logger.Crit("Failed packing storeCrossChainMessage message!")
			return syntheticTransactions

			// todo (@stefan) - return error
			// return nil, fmt.Errorf("failed packing submitOutOfNetworkMessage %w", err)
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
		data, err := MessageBusABI.Pack("notifyDeposit", transfer.Receiver, transfer.Amount)
		if err != nil {
			m.logger.Crit("Failed packing notifyDeposit message!")
			return syntheticTransactions
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

	return syntheticTransactions
}
