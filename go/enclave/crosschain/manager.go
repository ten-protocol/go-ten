package crosschain

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

const (
	ownerKeyHex = "6e384a07a01263518a18a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
)

type Manager struct {
	txOpts       *bind.TransactOpts
	contractABI  abi.ABI
	l1MessageBus gethcommon.Address
	l2MessageBus gethcommon.Address
	storage      db.Storage
	logger       gethlog.Logger
}

type CrossChainManager interface {
	GenerateMessageBusDeployTx() *types.Transaction
	//	SetL2MessageBusAddress(addr *gethcommon.Address)
	ProcessSyntheticTransactions(block *types.Block, receipts types.Receipts) error
	GetSyntheticTransactions(block *types.Block, receipts types.Receipts) types.Transactions
	GetSyntheticTransactionsBetween(fromBlock *types.Block, toBlock *types.Block, rollupState *state.StateDB) types.Transactions
	ExtractMessagesFromReceipts(receipts types.Receipts) []MessageBus.StructsCrossChainMessage
	ExtractMessagesFromReceipt(receipt *types.Receipt) []MessageBus.StructsCrossChainMessage
	GetOwner() *gethcommon.Address
}

func New(
	l1BusAddress *gethcommon.Address,
	l2BusAddress *gethcommon.Address,
	storage db.Storage, /*key *ecdsa.PrivateKey,*/
	chainId *big.Int,
	logger gethlog.Logger,
) CrossChainManager {
	contractABI, err := abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	if err != nil {
		panic(err) //panic?
	}

	key, _ := crypto.HexToECDSA(ownerKeyHex)
	txOpts, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		//log error todo::
		return nil
	}

	logger.Info(fmt.Sprintf("[CrossChain] L2 Cross Chain Owner Address: %s", txOpts.From.Hex()))

	//Key is derived, address is predictable, thus address of contract is predictible across all enclaves
	l2MessageBus := crypto.CreateAddress(txOpts.From, 0)

	//Start from 1 since 0 tx deploys system contract
	txOpts.Nonce = big.NewInt(1)

	return &Manager{
		l1MessageBus: *l1BusAddress,
		l2MessageBus: l2MessageBus,
		contractABI:  contractABI,
		txOpts:       txOpts,
		storage:      storage,
		logger:       logger,
	}
}

func (m *Manager) GetOwner() *gethcommon.Address {
	return &m.txOpts.From
}

func (m *Manager) GenerateMessageBusDeployTx() *types.Transaction {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0, //this should be fixed probably :/
		Value:    gethcommon.Big0,
		Gas:      5_000_000,       //requires above 1m gas to deploy wtf.
		GasPrice: gethcommon.Big0, //Synthetic transactions are on the house. Or the house.
		Data:     gethcommon.FromHex(MessageBus.MessageBusMetaData.Bin),
		To:       nil, //Geth requires nil instead of gethcommon.Address{} which equates to zero address in order to return receipt.
	})

	stx, err := m.txOpts.Signer(m.txOpts.From, tx)
	if err != nil {
		panic(err)
	}

	m.logger.Info(fmt.Sprintf("[CrossChain] Generated synthetic deployment transaction for the MessageBus contract %s - TX HASH: %s", m.l2MessageBus.Hex(), stx.Hash().Hex()))

	return stx
}

func (m *Manager) ProcessSyntheticTransactions(block *types.Block, receipts types.Receipts) error {
	if len(receipts) > 0 {
		m.lazilyLogReceiptChecksum(fmt.Sprintf("[CrossChain] Processing block: %s receipts: %d", block.Hash().Hex(), len(receipts)), receipts)
	}

	transactions := m.GetSyntheticTransactions(block, receipts)
	if len(transactions) > 0 {
		m.logger.Info(fmt.Sprintf("[CrossChain] Storing %d transactions for block %s", len(transactions), block.Hash().Hex()))
		m.lazilyLogChecksum("[CrossChain] Process synthetic transaction checksum", transactions)
		m.storage.StoreSyntheticTransactions(block.Hash(), transactions)
	}
	return nil
}

func (m *Manager) lazilyLogReceiptChecksum(msg string, receipts types.Receipts) {
	m.logger.Trace(msg, "Hash",
		gethlog.Lazy{Fn: func() string {
			hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
			hasher.Reset()
			for _, receipt := range receipts {
				var buffer bytes.Buffer
				receipt.EncodeRLP(&buffer)
				hasher.Write(buffer.Bytes())
			}
			var hash gethcommon.Hash
			hasher.Read(hash[:])
			return hash.Hex()
		}})
}

func (m *Manager) lazilyLogChecksum(msg string, transactions types.Transactions) {
	m.logger.Trace(msg, "Hash",
		gethlog.Lazy{Fn: func() string {
			hasher := sha3.NewLegacyKeccak256().(crypto.KeccakState)
			hasher.Reset()
			for _, tx := range transactions {
				var buffer bytes.Buffer
				tx.EncodeRLP(&buffer)
				hasher.Write(buffer.Bytes())
			}
			var hash gethcommon.Hash
			hasher.Read(hash[:])
			return hash.Hex()
		}})
}

func (m *Manager) GetSyntheticTransactionsBetween(fromBlock *types.Block, toBlock *types.Block, rollupState *state.StateDB) types.Transactions {
	transactions := make(types.Transactions, 0)

	//todo:: replace this with an iterator
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

		m.logger.Info(fmt.Sprintf("[CrossChain] Looking for transactions at block %s", b.Hash().Hex()))
		syntheticTransactions := m.storage.ReadSyntheticTransactions(b.Hash())
		transactions = append(transactions, syntheticTransactions...) //Ordering here might work in POBI, but might be weird for fast finality

		if b.NumberU64() < height {
			m.logger.Crit("block height is less than genesis height")
		}
		p, f := m.storage.ParentBlock(b)
		if !f {
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}
	m.lazilyLogChecksum("[CrossChain] Read synthetic transactions checksum", transactions)

	//Todo:: iteration order is reversed! This might cause unintended consequences!
	//Sign transactions and put proper nonces.
	startingNonce := rollupState.GetNonce(m.txOpts.From)

	signedTransactions := make(types.Transactions, 0)
	for idx, unsignedTransaction := range transactions {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    startingNonce + uint64(idx), //this should be fixed probably :/
			Value:    gethcommon.Big0,
			Gas:      5_000_000,
			GasPrice: gethcommon.Big0, //Synthetic transactions are on the house. Or the house.
			Data:     unsignedTransaction.Data(),
			To:       &m.l2MessageBus,
		})

		stx, err := m.txOpts.Signer(m.txOpts.From, tx)
		if err != nil {
			panic(err)
		}
		signedTransactions = append(signedTransactions, stx)
	}

	return signedTransactions
}

func (m *Manager) GetSyntheticTransactions(block *types.Block, receipts types.Receipts) types.Transactions {
	transactions := make(types.Transactions, 0)

	if len(receipts) == 0 {
		return transactions
	}

	/*if !VerifyReceiptHash(block, receipts) {
		m.logger.Crit("Receipts mismatch!")
		return transactions
	}*/

	messages := m.ExtractMessagesFromReceipts(receipts)

	m.logger.Info(fmt.Sprintf("[CrossChain] Found %d cross chain messages that will be submitted to L2!", len(messages)),
		"Block", block.Hash().Hex())

	for _, message := range messages {
		validAfter := big.NewInt(1)
		data, err := m.contractABI.Pack("submitOutOfNetworkMessage", &message, validAfter)
		if err != nil {
			panic(err)
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    0, //This gets replaced later on
			Value:    gethcommon.Big0,
			Gas:      5_000_000,
			GasPrice: gethcommon.Big0, //Synthetic transactions are on the house. Or the house.
			Data:     data,
			To:       &m.l2MessageBus,
		})

		m.logger.Info(fmt.Sprintf("[CrossChain] Creating synthetic tx for cross chain message to L2. From: %s Topic: %s",
			message.Sender.Hex(),
			fmt.Sprint(message.Topic)), "Block", block.Hash().Hex())

		transactions = append(transactions, tx)
	}

	return transactions
}

func (m *Manager) ExtractMessagesFromReceipts(receipts types.Receipts) []MessageBus.StructsCrossChainMessage {
	messages := make([]MessageBus.StructsCrossChainMessage, 0)

	for _, receipt := range receipts {
		extractedCrossChainMessages := m.ExtractMessagesFromReceipt(receipt)
		messages = append(messages, extractedCrossChainMessages...)
	}

	return messages
}

func (m *Manager) ExtractMessagesFromReceipt(receipt *types.Receipt) []MessageBus.StructsCrossChainMessage {
	if receiptMightContainPublishedMessage(receipt) {
		events := m.extractPublishedMessages(receipt)
		return convertToMessages(events)
	}

	return make([]MessageBus.StructsCrossChainMessage, 0)
}

func (m *Manager) extractPublishedMessages(receipt *types.Receipt) []MessageBus.MessageBusLogMessagePublished {
	events := make([]MessageBus.MessageBusLogMessagePublished, 0)

	m.logger.Info(fmt.Sprintf("[CrossChain] Extracting %d logs from receipt for %s", len(receipt.Logs), receipt.TxHash.Hex()))

	for _, log := range receipt.Logs {
		//event, err := m.l2MessageBus.ParseLogMessagePublished(*log)
		event := m.extractPublishedMessage(log)
		if event != nil {
			events = append(events, *event)
		}
	}

	return events
}

func (m *Manager) extractPublishedMessage(log *types.Log) *MessageBus.MessageBusLogMessagePublished {
	eventABI := m.contractABI.Events["LogMessagePublished"]
	//Unpack only relevant logs.
	if log.Topics[0] != eventABI.ID {
		return nil
	}

	//Unpack only from our system contracts.
	//Otherwise someone can just post a clone contract that matches event sig and token goes puf.
	//todo:: perhaps dont convert to hex everytime
	if (log.Address.Hex() != m.l1MessageBus.Hex()) && (log.Address.Hex() != m.l2MessageBus.Hex()) {
		return nil
	}

	m.logger.Info(fmt.Sprintf("[CrossChain] Event from message bus %s found!", log.Address.Hex()))

	var event MessageBus.MessageBusLogMessagePublished
	m.contractABI.UnpackIntoInterface(&event, "LogMessagePublished", log.Data)

	m.logger.Trace(fmt.Sprintf("[CrossChain] Event extracted - %+v", event))
	return &event
}

func VerifyReceiptHash(block *types.Block, receipts types.Receipts) bool {
	hash := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	return block.ReceiptHash().Hex() == hash.Hex()
}

func receiptMightContainPublishedMessage(receipt *types.Receipt) bool {
	//todo:: check bloom filter of receipt after figuring out how :|
	return true
}

func convertToMessages(events []MessageBus.MessageBusLogMessagePublished) []MessageBus.StructsCrossChainMessage {
	messages := make([]MessageBus.StructsCrossChainMessage, 0)

	for _, event := range events {
		msg := createCrossChainMessage(event)
		messages = append(messages, msg)
	}

	return messages
}

func createCrossChainMessage(event MessageBus.MessageBusLogMessagePublished) MessageBus.StructsCrossChainMessage {
	return MessageBus.StructsCrossChainMessage{
		Sender:   event.Sender,
		Sequence: event.Sequence,
		Nonce:    event.Nonce,
		Topic:    event.Topic,
		Payload:  event.Payload,
	}
}
