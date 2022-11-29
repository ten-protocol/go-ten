package crosschain

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

const (
	ownerKeyHex = "6e384a07a01263518a18a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
)

type obscuroMessageBusManager struct {
	messageBusAddress *gethcommon.Address
	storage           db.Storage
	logger            gethlog.Logger
	wallet            wallet.Wallet
}

func NewObscuroMessageBusManager(
	storage db.Storage, /*key *ecdsa.PrivateKey,*/
	chainID *big.Int,
	logger gethlog.Logger,
) ObscuroCrossChainManager {
	key, _ := crypto.HexToECDSA(ownerKeyHex)
	wallet := wallet.NewInMemoryWalletFromPK(chainID, key, logger)

	logger.Info(fmt.Sprintf("[CrossChain] L2 Cross Chain Owner Address: %s", wallet.Address().Hex()))

	// Key is derived, address is predictable, thus address of contract is predictable across all enclaves
	l2MessageBus := crypto.CreateAddress(wallet.Address(), 0)

	return &obscuroMessageBusManager{
		messageBusAddress: &l2MessageBus,
		storage:           storage,
		logger:            logger,
		wallet:            wallet,
	}
}

func (m *obscuroMessageBusManager) IsSyntheticTransaction(transaction common.L2Tx) bool {
	sender, err := rpc.GetSender(&transaction)
	if err != nil {
		return false
	}

	return bytes.Equal(sender.Bytes(), m.GetOwner().Bytes())
}

// GetOwner - Returns the address of the identity owning the message bus.
func (m *obscuroMessageBusManager) GetOwner() common.L2Address {
	return m.wallet.Address()
}

// GetBusAddress - Returns the L2 address of the message bus contract.
// TODO: Figure out how to expose the deployed contract to the external world. Perhaps extract event from contract construction?
func (m *obscuroMessageBusManager) GetBusAddress() *common.L2Address {
	return m.messageBusAddress
}

// DeriveOwner - Generates the key pair that will be used to transact with the L2 message bus.
func (m *obscuroMessageBusManager) DeriveOwner(seed []byte) (*common.L2Address, error) {
	// TODO: Implement with cryptography epic!
	return m.messageBusAddress, nil
}

// GenerateMessageBusDeployTx - Returns a signed message bus deployment transaction.
func (m *obscuroMessageBusManager) GenerateMessageBusDeployTx() (*common.L2Tx, error) {
	tx := &types.LegacyTx{
		Nonce:    0, // this should be fixed probably :/
		Value:    gethcommon.Big0,
		Gas:      5_000_000,       // requires above 1m gas to deploy wtf.
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     gethcommon.FromHex(MessageBus.MessageBusMetaData.Bin),
		To:       nil, // Geth requires nil instead of gethcommon.Address{} which equates to zero address in order to return receipt.
	}

	stx, err := m.wallet.SignTransaction(tx)
	if err != nil {
		return nil, err
	}

	m.logger.Trace(fmt.Sprintf("[CrossChain] Generated synthetic deployment transaction for the MessageBus contract %s - TX HASH: %s", m.messageBusAddress.Hex(), stx.Hash().Hex()))

	return stx, nil
}

// ExtractLocalMessages - Finds relevant logs in the receipts and converts them to cross chain messages.
func (m *obscuroMessageBusManager) ExtractLocalMessages(receipts common.L2Receipts) (common.CrossChainMessages, error) {
	logs, err := filterLogsFromReceipts(receipts, m.messageBusAddress, &CrossChainEventID)
	if err != nil {
		m.logger.Error("[CrossChain] Error extracting logs from L2 message bus!", "Error", err)
		return make(common.CrossChainMessages, 0), err
	}

	messages, err := convertLogsToMessages(logs, CrossChainEventName, ContractABI)
	if err != nil {
		m.logger.Error("[CrossChain] Error converting messages from L2 message bus!", "Error", err)
		return make(common.CrossChainMessages, 0), err
	}

	return messages, nil
}

// SubmitRemoteMessagesLocally - Submits messages saved between the from and to blocks on chain using the provided function bindings.
func (m *obscuroMessageBusManager) SubmitRemoteMessagesLocally(
	fromBlock *common.L1Block,
	toBlock *common.L1Block,
	rollupState *state.StateDB,
	processTxCall OnChainEVMExecutorFunc,
	processOffChainMessage OffChainEVMCallFunc,
) error {
	transactions := m.retrieveSyntheticTransactionsBetween(fromBlock, toBlock, rollupState)
	m.logger.Trace("[CrossChain] Retrieved synthetic transactions",
		"Count", len(transactions),
		"FromBlock", common.ShortHash(fromBlock.Hash()),
		"ToBlock", common.ShortHash(toBlock.Hash()))

	if len(transactions) > 0 {
		syntheticTransactionsResponses := processTxCall(transactions)
		synthReceipts := make([]*types.Receipt, len(syntheticTransactionsResponses))
		if len(syntheticTransactionsResponses) != len(transactions) {
			m.logger.Crit("Sanity check. Some synthetic transactions failed.")
			return errors.New("evm failed to generate responses for every transaction")
		}

		i := 0
		for _, resp := range syntheticTransactionsResponses {
			rec, ok := resp.(*types.Receipt)
			if !ok { // Еxtract reason for failing deposit.
				// TODO - Handle the case of an error (e.g. insufficient funds).
				m.logger.Crit("Sanity check. Expected a receipt", log.ErrKey, resp)
				return errors.New("receipt missing for a guaranteed synthetic transaction")
			}

			if rec.Status == 0 {
				failingTx := transactions[i]
				txCallMessage := types.NewMessage(
					m.GetOwner(),
					failingTx.To(),
					rollupState.GetNonce(m.GetOwner()),
					failingTx.Value(),
					failingTx.Gas(),
					gethcommon.Big0,
					gethcommon.Big0,
					gethcommon.Big0,
					failingTx.Data(),
					failingTx.AccessList(),
					false)

				res, err := processOffChainMessage(txCallMessage)
				m.logger.Crit("Synthetic transaction failed!", log.ErrKey, err, "result", res)
				return fmt.Errorf("synthetic transaction failed. error: %w result: %+v", err, res)
			}

			synthReceipts[i] = rec
			i++
		}
	}

	return nil
}

// retrieveSyntheticTransactionsBetween - Iterates the blocks backwards and returns the synthetic transaction
// TODO: Fix ordering of transactions, currently it is irrelevant.
// TODO: Do not extract messages below their consistency level. Irrelevant security wise.
// TODO: Surface errors
func (m *obscuroMessageBusManager) retrieveSyntheticTransactionsBetween(fromBlock *common.L1Block, toBlock *common.L1Block, rollupState *state.StateDB) common.L2Transactions {
	messages := make(common.CrossChainMessages, 0)

	from := common.GenesisBlock.Hash()
	height := common.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !m.storage.IsAncestor(toBlock, fromBlock) {
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}
	// Iterate through the blocks.
	b := toBlock
	for {
		if bytes.Equal(b.Hash().Bytes(), from.Bytes()) {
			break
		}

		m.logger.Trace(fmt.Sprintf("[CrossChain] Looking for transactions at block %s", b.Hash().Hex()))
		messagesForBlock := m.storage.ReadL1Messages(b.Hash())
		messages = append(messages, messagesForBlock...) // Ordering here might work in POBI, but might be weird for fast finality

		// No deposits before genesis.
		if b.NumberU64() < height {
			m.logger.Crit("block height is less than genesis height")
		}
		p, err := m.storage.ParentBlock(b)
		if err != nil {
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}

	// Get current nonce for this stateDB.
	// There can be forks thus we cannot trust the wallet.
	startingNonce := rollupState.GetNonce(m.GetOwner())

	signedTransactions := make(types.Transactions, 0)
	for idx, message := range messages {
		delayInBlocks := big.NewInt(int64(message.ConsistencyLevel))
		data, err := ContractABI.Pack("submitOutOfNetworkMessage", message, delayInBlocks)
		if err != nil {
			m.logger.Crit("[CrossChain] Failed packing submitOutOfNetwork message!")
			return signedTransactions

			// TODO: return error
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

		stx, err := m.wallet.SignTransaction(tx)
		if err != nil {
			panic(err)
		}
		signedTransactions = append(signedTransactions, stx)
	}

	return signedTransactions
}
