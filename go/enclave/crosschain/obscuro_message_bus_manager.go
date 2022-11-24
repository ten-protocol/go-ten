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

	// Key is derived, address is predictable, thus address of contract is predictible across all enclaves
	l2MessageBus := crypto.CreateAddress(wallet.Address(), 0)

	return &obscuroMessageBusManager{
		messageBusAddress: &l2MessageBus,
		storage:           storage,
		logger:            logger,
		wallet:            wallet,
	}
}

func (m *obscuroMessageBusManager) GetOwner() common.L2Address {
	return m.wallet.Address()
}

func (m *obscuroMessageBusManager) GetBusAddress() *common.L2Address {
	return m.messageBusAddress
}

func (m *obscuroMessageBusManager) DeriveOwner(seed []byte) (*common.L2Address, error) {
	// TODO: Implement with cryptography epic!
	return m.messageBusAddress, nil
}

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
				return fmt.Errorf("synthetic transaction failed. error: %+v result: %+v", err, res)
			}

			synthReceipts[i] = rec
			i++
		}
	}

	return nil
}

func (m *obscuroMessageBusManager) retrieveSyntheticTransactionsBetween(fromBlock *common.L1Block, toBlock *common.L1Block, rollupState *state.StateDB) common.L2Transactions {
	transactions := make(common.L2Transactions, 0)

	from := common.GenesisBlock.Hash()
	height := common.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !m.storage.IsAncestor(toBlock, fromBlock) {
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}

	b := toBlock
	for {
		if bytes.Equal(b.Hash().Bytes(), from.Bytes()) {
			break
		}

		m.logger.Trace(fmt.Sprintf("[CrossChain] Looking for transactions at block %s", b.Hash().Hex()))
		syntheticTransactions := m.storage.ReadSyntheticTransactions(b.Hash())
		transactions = append(transactions, syntheticTransactions...) // Ordering here might work in POBI, but might be weird for fast finality

		if b.NumberU64() < height {
			m.logger.Crit("block height is less than genesis height")
		}
		p, f := m.storage.ParentBlock(b)
		if !f {
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}
	lazilyLogChecksum("[CrossChain] Read synthetic transactions checksum", transactions, m.logger)

	// Todo:: iteration order is reversed! This might cause unintended consequences!
	// Sign transactions and put proper nonces.
	startingNonce := rollupState.GetNonce(m.GetOwner())

	signedTransactions := make(types.Transactions, 0)
	for idx, unsignedTransaction := range transactions {
		tx := &types.LegacyTx{
			Nonce:    startingNonce + uint64(idx), // this should be fixed probably :/
			Value:    gethcommon.Big0,
			Gas:      5_000_000,
			GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
			Data:     unsignedTransaction.Data(),
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
