package crosschain

import (
	"bytes"
	"math/big"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

const (
	ownerKeyHex = "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
)

type Manager struct {
	txOpts       *bind.TransactOpts
	contractABI  abi.ABI
	l1MessageBus *gethcommon.Address
	l2MessageBus *gethcommon.Address
	storage      db.Storage
	logger       gethlog.Logger
}

type CrossChainManager interface {
	ProcessSyntheticTransactions(block *types.Block, receipts []*types.Receipt) error
	GetSyntheticTransactions(block *types.Block, receipts []*types.Receipt) types.Transactions
	GetSyntheticTransactionsBetween(fromBlock *types.Block, toBlock *types.Block) types.Transactions
	ExtractMessagesFromReceipts(receipts []*types.Receipt) []*MessageBus.StructsCrossChainMessage
	ExtractMessagesFromReceipt(receipt *types.Receipt) []*MessageBus.StructsCrossChainMessage
}

func New(
	l1BusAddress *gethcommon.Address,
	l2BusAddress *gethcommon.Address,
	storage db.Storage, /*key *ecdsa.PrivateKey,*/
	chainId *big.Int,
	logger gethlog.Logger,
) *Manager {
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

	return &Manager{
		l1MessageBus: l1BusAddress,
		l2MessageBus: l2BusAddress,
		contractABI:  contractABI,
		txOpts:       txOpts,
		storage:      storage,
		logger:       logger,
	}
}

func (m *Manager) ProcessSyntheticTransactions(block *types.Block, receipts []*types.Receipt) error {

	transactions := m.GetSyntheticTransactions(block, receipts)
	m.storage.StoreSyntheticTransactions(block.Hash(), transactions)

	return nil
}

func (m *Manager) GetSyntheticTransactionsBetween(fromBlock *types.Block, toBlock *types.Block) types.Transactions {
	transactions := make(types.Transactions, 1)

	//todo:: replace this with an iterator
	from := common.GenesisBlock.Hash()
	height := common.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !m.storage.IsAncestor(toBlock, fromBlock) {
			//todo:: logger
			m.logger.Crit("Synthetic transactions can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}

	b := toBlock
	for {
		if bytes.Equal(b.Hash().Bytes(), from.Bytes()) {
			break
		}

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

	return transactions
}

func (m *Manager) GetSyntheticTransactions(block *types.Block, receipts []*types.Receipt) types.Transactions {
	transactions := make(types.Transactions, 1)

	if !VerifyReceiptHash(block, receipts) {
		return transactions
	}

	messages := m.ExtractMessagesFromReceipts(receipts)

	for idx, message := range messages {
		data, err := m.contractABI.Pack("submitOutOfNetworkMessage", *message, big.NewInt(0))
		if err != nil {
			panic(err)
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    m.txOpts.Nonce.Uint64() + uint64(idx), //this should be fixed probably :/
			Value:    gethcommon.Big0,
			Gas:      1_000_000,
			GasPrice: gethcommon.Big0, //Synthetic transactions are on the house. Or the house.
			Data:     data,
			To:       m.l2MessageBus,
		})

		stx, err := m.txOpts.Signer(m.txOpts.From, tx)
		if err != nil {
			panic(err)
		}
		transactions = append(transactions, stx)
	}

	return transactions
}

func (m *Manager) ExtractMessagesFromReceipts(receipts []*types.Receipt) []*MessageBus.StructsCrossChainMessage {
	messages := make([]*MessageBus.StructsCrossChainMessage, 10)

	for _, receipt := range receipts {
		messages = append(messages, m.ExtractMessagesFromReceipt(receipt)...)
	}

	return messages
}

func (m *Manager) ExtractMessagesFromReceipt(receipt *types.Receipt) []*MessageBus.StructsCrossChainMessage {
	if receiptMightContainPublishedMessage(receipt) {
		events := m.extractPublishedMessages(receipt)
		return convertToMessages(events)
	}

	return make([]*MessageBus.StructsCrossChainMessage, 0)
}

func (m *Manager) extractPublishedMessages(receipt *types.Receipt) []*MessageBus.MessageBusLogMessagePublished {
	events := make([]*MessageBus.MessageBusLogMessagePublished, len(receipt.Logs))

	for _, log := range receipt.Logs {
		//event, err := m.l2MessageBus.ParseLogMessagePublished(*log)
		event := m.extractPublishedMessage(log)
		if event != nil {
			continue
		}

		events = append(events, event)
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
	if log.Address.Hex() != m.l1MessageBus.Hex() || log.Address.Hex() != m.l2MessageBus.Hex() {
		return nil
	}

	var event MessageBus.MessageBusLogMessagePublished
	m.contractABI.UnpackIntoInterface(event, "LogMessagePublished", log.Data)
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

func convertToMessages(events []*MessageBus.MessageBusLogMessagePublished) []*MessageBus.StructsCrossChainMessage {
	messages := make([]*MessageBus.StructsCrossChainMessage, len(events))

	for _, event := range events {
		msg := createCrossChainMessage(event)
		messages = append(messages, &msg)
	}

	return messages
}

func createCrossChainMessage(event *MessageBus.MessageBusLogMessagePublished) MessageBus.StructsCrossChainMessage {
	return MessageBus.StructsCrossChainMessage{
		Sender:   event.Sender,
		Sequence: event.Sequence,
		Nonce:    event.Nonce,
		Topic:    event.Topic,
		Payload:  event.Payload,
	}
}
