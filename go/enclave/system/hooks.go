package system

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/TransactionPostProcessor"
	"github.com/ten-protocol/go-ten/contracts/generated/ZenBase"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/wallet"
)

var (
	transactionPostProcessorABI, _ = abi.JSON(strings.NewReader(TransactionPostProcessor.TransactionPostProcessorMetaData.ABI))
	ErrNoTransactions              = fmt.Errorf("no transactions")
)

type SystemContractCallbacks interface {
	GetOwner() gethcommon.Address
	Initialize(batch *core.Batch, receipts types.Receipts) error
	Load() error
	CreateOnBatchEndTransaction(ctx context.Context, stateDB *state.StateDB, transactions common.L2Transactions, receipts types.Receipts) (*types.Transaction, error)
	TransactionPostProcessor() *gethcommon.Address
	VerifyOnBlockReceipt(transactions common.L2Transactions, receipt *types.Receipt) (bool, error)
}

type systemContractCallbacks struct {
	transactionsAnalyzerAddress *gethcommon.Address
	ownerWallet                 wallet.Wallet
	storage                     storage.Storage

	logger gethlog.Logger
}

func NewSystemContractCallbacks(ownerWallet wallet.Wallet, logger gethlog.Logger) SystemContractCallbacks {
	return &systemContractCallbacks{
		transactionsAnalyzerAddress: nil,
		ownerWallet:                 ownerWallet,
		logger:                      logger,
		storage:                     nil,
	}
}

func (s *systemContractCallbacks) TransactionPostProcessor() *gethcommon.Address {
	return s.transactionsAnalyzerAddress
}

func (s *systemContractCallbacks) GetOwner() gethcommon.Address {
	return s.ownerWallet.Address()
}

func (s *systemContractCallbacks) Load() error {
	s.logger.Info("Load: Initializing system contracts")

	if s.storage == nil {
		s.logger.Error("Load: Storage is not set")
		return fmt.Errorf("storage is not set")
	}

	batchSeqNo := uint64(1)
	s.logger.Debug("Load: Fetching batch", "batchSeqNo", batchSeqNo)
	batch, err := s.storage.FetchBatchBySeqNo(context.Background(), batchSeqNo)
	if err != nil {
		s.logger.Error("Load: Failed fetching batch", "batchSeqNo", batchSeqNo, "error", err)
		return fmt.Errorf("failed fetching batch %w", err)
	}

	if len(batch.Transactions) < 2 {
		s.logger.Error("Load: Genesis batch does not have enough transactions", "batchSeqNo", batchSeqNo, "transactionCount", len(batch.Transactions))
		return fmt.Errorf("genesis batch does not have enough transactions")
	}

	receipt, err := s.storage.GetTransactionReceipt(context.Background(), batch.Transactions[1].Hash())
	if err != nil {
		s.logger.Error("Load: Failed fetching receipt", "transactionHash", batch.Transactions[1].Hash().Hex(), "error", err)
		return fmt.Errorf("failed fetching receipt %w", err)
	}

	addresses, err := DeriveAddresses(receipt)
	if err != nil {
		s.logger.Error("Load: Failed deriving addresses", "error", err, "receiptHash", receipt.TxHash.Hex())
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) initializeRequiredAddresses(addresses SystemContractAddresses) error {
	if addresses["TransactionsAnalyzer"] == nil {
		return fmt.Errorf("required contract address TransactionsAnalyzer is nil")
	}

	s.transactionsAnalyzerAddress = addresses["TransactionsAnalyzer"]

	return nil
}

func (s *systemContractCallbacks) Initialize(batch *core.Batch, receipts types.Receipts) error {
	s.logger.Info("Initialize: Starting initialization of system contracts", "batchSeqNo", batch.SeqNo)

	if len(receipts) < 2 {
		s.logger.Error("Initialize: Genesis batch does not have enough receipts", "expected", 2, "got", len(receipts))
		return fmt.Errorf("genesis batch does not have enough receipts")
	}

	receiptIndex := 1
	s.logger.Debug("Initialize: Deriving addresses from receipt", "receiptIndex", receiptIndex, "transactionHash", receipts[receiptIndex].TxHash.Hex())
	addresses, err := DeriveAddresses(receipts[receiptIndex])
	if err != nil {
		s.logger.Error("Initialize: Failed deriving addresses", "error", err, "receiptHash", receipts[receiptIndex].TxHash.Hex())
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	s.logger.Info("Initialize: Initializing required addresses", "addresses", addresses)
	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) CreateOnBatchEndTransaction(ctx context.Context, l2State *state.StateDB, transactions common.L2Transactions, receipts types.Receipts) (*types.Transaction, error) {
	if s.transactionsAnalyzerAddress == nil {
		s.logger.Debug("CreateOnBatchEndTransaction: TransactionsAnalyzerAddress is nil, skipping transaction creation")
		return nil, nil
	}

	if len(transactions) == 0 {
		s.logger.Debug("CreateOnBatchEndTransaction: Batch has no transactions, skipping transaction creation")
		return nil, ErrNoTransactions
	}

	nonceForSyntheticTx := l2State.GetNonce(s.GetOwner())
	s.logger.Debug("CreateOnBatchEndTransaction: Retrieved nonce for synthetic transaction", "nonce", nonceForSyntheticTx)

	solidityTransactions := make([]TransactionPostProcessor.StructsTransaction, 0)

	type statusWithGasUsed struct {
		status  bool
		gasUsed uint64
	}

	txSuccessMap := map[gethcommon.Hash]statusWithGasUsed{}
	for _, receipt := range receipts {
		txSuccessMap[receipt.TxHash] = statusWithGasUsed{
			status:  receipt.Status == types.ReceiptStatusSuccessful,
			gasUsed: receipt.GasUsed,
		}
	}

	for _, tx := range transactions {
		// Start of Selection
		transaction := TransactionPostProcessor.StructsTransaction{
			Nonce:      big.NewInt(int64(tx.Nonce())),
			GasPrice:   tx.GasPrice(),
			GasLimit:   big.NewInt(int64(tx.Gas())),
			Value:      tx.Value(),
			Data:       tx.Data(),
			Successful: txSuccessMap[tx.Hash()].status,
			GasUsed:    txSuccessMap[tx.Hash()].gasUsed,
		}
		if tx.To() != nil {
			transaction.To = *tx.To()
		} else {
			transaction.To = gethcommon.Address{} // Zero address - contract deployment
		}

		sender, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			s.logger.Error("CreateOnBatchEndTransaction: Failed to recover sender address", "error", err, "transactionHash", tx.Hash().Hex())
			return nil, fmt.Errorf("failed to recover sender address: %w", err)
		}
		transaction.From = sender

		solidityTransactions = append(solidityTransactions, transaction)
		s.logger.Debug("CreateOnBatchEndTransaction: Encoded transaction", "transactionHash", tx.Hash().Hex(), "sender", sender.Hex())
	}

	data, err := transactionPostProcessorABI.Pack("onBlock", solidityTransactions)
	if err != nil {
		s.logger.Error("CreateOnBatchEndTransaction: Failed packing onBlock data", "error", err)
		return nil, fmt.Errorf("failed packing onBlock() %w", err)
	}

	tx := &types.LegacyTx{
		Nonce:    nonceForSyntheticTx,
		Value:    gethcommon.Big0,
		Gas:      500_000_000,
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     data,
		To:       s.transactionsAnalyzerAddress,
	}

	s.logger.Debug("CreateOnBatchEndTransaction: Signing transaction", "to", s.transactionsAnalyzerAddress.Hex(), "nonce", nonceForSyntheticTx)
	signedTx, err := s.ownerWallet.SignTransaction(tx)
	if err != nil {
		s.logger.Error("CreateOnBatchEndTransaction: Failed signing transaction", "error", err)
		return nil, fmt.Errorf("failed signing transaction %w", err)
	}

	s.logger.Info("CreateOnBatchEndTransaction: Successfully created signed transaction", "transactionHash", signedTx.Hash().Hex())
	return signedTx, nil
}

func (s *systemContractCallbacks) VerifyOnBlockReceipt(transactions common.L2Transactions, receipt *types.Receipt) (bool, error) {
	if receipt.Status != types.ReceiptStatusSuccessful {
		s.logger.Error("VerifyOnBlockReceipt: Transaction failed", "transactionHash", receipt.TxHash.Hex())
		return false, fmt.Errorf("transaction failed")
	}

	if len(receipt.Logs) == 0 {
		s.logger.Error("VerifyOnBlockReceipt: Transaction has no logs", "transactionHash", receipt.TxHash.Hex())
		return false, fmt.Errorf("transaction has no logs")
	}

	abi, err := ZenBase.ZenBaseMetaData.GetAbi()
	if err != nil {
		s.logger.Error("VerifyOnBlockReceipt: Failed to get ABI", "error", err)
		return false, fmt.Errorf("failed to get ABI %w", err)
	}

	if len(receipt.Logs) == 0 {
		s.logger.Error("VerifyOnBlockReceipt: Synthetic transaction has no logs", "transactionHash", receipt.TxHash.Hex())
		return false, fmt.Errorf("no logs in onBlockReceipt")
	}

	// Find the TransactionsConverted event in the onBlockReceipt and verify the number of transactions converted
	// matches the number of transactions in the batch. Mostly paranoia code.
	for _, log := range receipt.Logs {
		if len(log.Topics) > 0 && log.Topics[0] == abi.Events["TransactionsConverted"].ID { // TransactionsConverted event signature
			if len(log.Data) != 32 {
				s.logger.Error("VerifyOnBlockReceipt: Invalid data length for TransactionsConverted event", "expected", 32, "got", len(log.Data))
				return false, fmt.Errorf("invalid data length for TransactionsConverted event")
			}
			transactionsConverted := new(big.Int).SetBytes(log.Data)
			if transactionsConverted.Uint64() != uint64(len(transactions)) {
				s.logger.Error("VerifyOnBlockReceipt: Mismatch in TransactionsConverted event", "expected", len(transactions), "got", transactionsConverted.Uint64())
				return false, fmt.Errorf("mismatch in TransactionsConverted event: expected %d, got %d", len(transactions), transactionsConverted.Uint64())
			}
			break
		}
	}

	for _, log := range receipt.Logs {
		if log.Topics[0] != abi.Events["TransactionProcessed"].ID {
			continue
		}

		decodedLog, err := abi.Unpack("TransactionProcessed", log.Data)
		if err != nil {
			s.logger.Error("VerifyOnBlockReceipt: Failed to unpack log", "error", err, "log", log)
			return false, fmt.Errorf("failed to unpack log %w", err)
		}
		s.logger.Debug("VerifyOnBlockReceipt: Decoded log", "log", decodedLog)
	}

	s.logger.Debug("VerifyOnBlockReceipt: Transaction successful", "transactionHash", receipt.TxHash.Hex())
	return true, nil
}
