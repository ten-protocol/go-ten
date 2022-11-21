package crosschain

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/contracts/messagebuscontract/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type mainNetExtractor struct {
	busAddress   *common.L1Address
	l2MessageBus *common.L2Address //TODO: remove this
	storage      db.Storage
	logger       gethlog.Logger
	contractABI  abi.ABI
}

func NewMainNetExtractor(busAddress *common.L1Address, l2BusAddress *common.L2Address, storage db.Storage, logger gethlog.Logger) MainNetMessageExtractor {

	contractABI, err := abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	if err != nil {
		logger.Crit("Unable to initialize MainNetExtractor due to error when parsing contract ABI!")
	}

	return &mainNetExtractor{
		busAddress:   busAddress,
		l2MessageBus: l2BusAddress,
		storage:      storage,
		logger:       logger,
		contractABI:  contractABI,
	}
}

func (m *mainNetExtractor) ProcessCrossChainMessages(block *common.L1Block, receipts common.L1Receipts) error {
	areReceiptsValid := VerifyReceiptHash(block, receipts)

	if !areReceiptsValid {
		m.logger.Error("[CrossChain] Invalid receipts submitted", "block", common.ShortHash(block.Hash()))
		return fmt.Errorf("receipts do not match the receipt root for the block")
	}

	if len(receipts) == 0 {
		//Error if block receipts root does not match receipts hash
		//else nil
		return nil
	}

	lazilyLogReceiptChecksum(fmt.Sprintf("[CrossChain] Processing block: %s receipts: %d", block.Hash().Hex(), len(receipts)), receipts, m.logger)
	transactions, err := m.getSyntheticTransactions(block, receipts)

	if err != nil {
		return err
	}

	if len(transactions) > 0 {
		m.logger.Trace(fmt.Sprintf("[CrossChain] Storing %d transactions for block %s", len(transactions), block.Hash().Hex()))
		lazilyLogChecksum("[CrossChain] Process synthetic transaction checksum", transactions, m.logger)
		m.storage.StoreSyntheticTransactions(block.Hash(), transactions)
	}

	return nil
}

func (m *mainNetExtractor) GetBusAddress() *common.L1Address {
	return m.busAddress
}

func (m *mainNetExtractor) getSyntheticTransactions(block *common.L1Block, receipts common.L1Receipts) (common.L2Transactions, error) {
	transactions := make(common.L2Transactions, 0)

	if len(receipts) == 0 {
		return transactions, nil
	}

	eventId := m.contractABI.Events["LogMessagePublished"].ID
	logs, err := filterLogsFromReceipts(receipts, m.GetBusAddress(), &eventId)

	if err != nil {
		m.logger.Error("[CrossChain]", "Error", err)
		return transactions, err
	}
	m.logger.Trace("[CrossChain] extracted logs", "logCount", len(logs))

	messages, err := convertLogsToMessages(logs, "LogMessagePublished", m.contractABI)

	if err != nil {
		m.logger.Error("[CrossChain]", "Error", err)
		return transactions, err
	}

	m.logger.Trace(fmt.Sprintf("[CrossChain] Found %d cross chain messages that will be submitted to L2!", len(messages)),
		"Block", block.Hash().Hex())

	for _, message := range messages {
		validAfter := big.NewInt(1)
		data, err := m.contractABI.Pack("submitOutOfNetworkMessage", &message, validAfter)
		if err != nil {
			return transactions, fmt.Errorf("Failed packing submitOutOfNetworkMessage %+v", err)
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    0, //This gets replaced later on
			Value:    gethcommon.Big0,
			Gas:      5_000_000,
			GasPrice: gethcommon.Big0, //Synthetic transactions are on the house. Or the house.
			Data:     data,
			To:       m.l2MessageBus,
		})

		m.logger.Trace(fmt.Sprintf("[CrossChain] Creating synthetic tx for cross chain message to L2. From: %s Topic: %s",
			message.Sender.Hex(),
			fmt.Sprint(message.Topic)), "Block", block.Hash().Hex())

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
