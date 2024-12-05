package system

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/contracts/generated/PublicCallbacks"
	"github.com/ten-protocol/go-ten/contracts/generated/TransactionPostProcessor"
	"github.com/ten-protocol/go-ten/contracts/generated/ZenBase"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

var (
	transactionPostProcessorABI, _ = abi.JSON(strings.NewReader(TransactionPostProcessor.TransactionPostProcessorMetaData.ABI))
	publicCallbacksABI, _          = abi.JSON(strings.NewReader(PublicCallbacks.PublicCallbacksMetaData.ABI))
	ErrNoTransactions              = fmt.Errorf("no transactions")
)

type SystemContractCallbacks interface {
	// Getters
	PublicCallbackHandler() *gethcommon.Address
	TransactionPostProcessor() *gethcommon.Address
	SystemContractsUpgrader() *gethcommon.Address

	// Initialization
	Initialize(batch *core.Batch, receipts types.Receipt, msgBusManager SystemContractsInitializable) error
	Load() error

	// Usage
	CreateOnBatchEndTransaction(ctx context.Context, stateDB *state.StateDB, results core.TxExecResults) (*types.Transaction, error)
	CreatePublicCallbackHandlerTransaction(ctx context.Context, stateDB *state.StateDB) (*types.Transaction, error)
	VerifyOnBlockReceipt(transactions common.L2Transactions, receipt *types.Receipt) (bool, error)
}

type SystemContractsInitializable interface {
	Initialize(SystemContractAddresses) error
}

type systemContractCallbacks struct {
	transactionsPostProcessorAddress *gethcommon.Address
	storage                          storage.Storage
	systemAddresses                  SystemContractAddresses
	systemContractsUpgrader          *gethcommon.Address

	logger gethlog.Logger
}

func NewSystemContractCallbacks(storage storage.Storage, upgrader *gethcommon.Address, logger gethlog.Logger) SystemContractCallbacks {
	return &systemContractCallbacks{
		transactionsPostProcessorAddress: nil,
		logger:                           logger,
		storage:                          storage,
		systemAddresses:                  make(SystemContractAddresses),
		systemContractsUpgrader:          upgrader,
	}
}

func (s *systemContractCallbacks) SystemContractsUpgrader() *gethcommon.Address {
	return s.systemContractsUpgrader
}

func (s *systemContractCallbacks) TransactionPostProcessor() *gethcommon.Address {
	return s.transactionsPostProcessorAddress
}

func (s *systemContractCallbacks) PublicCallbackHandler() *gethcommon.Address {
	return s.systemAddresses["PublicCallbacks"]
}

func (s *systemContractCallbacks) Load() error {
	s.logger.Info("Load: Initializing system contracts")

	if s.storage == nil {
		s.logger.Error("Load: Storage is not set")
		return fmt.Errorf("storage is not set")
	}

	batchSeqNo := uint64(2)
	s.logger.Debug("Load: Fetching batch", "batchSeqNo", batchSeqNo)
	batch, err := s.storage.FetchBatchBySeqNo(context.Background(), batchSeqNo)
	if err != nil {
		s.logger.Error("Load: Failed fetching batch", "batchSeqNo", batchSeqNo, "error", err)
		return fmt.Errorf("failed fetching batch %w", err)
	}

	tx, err := SystemDeployerInitTransaction(s.logger, *s.systemContractsUpgrader)
	if err != nil {
		s.logger.Error("Load: Failed creating system deployer init transaction", "error", err)
		return fmt.Errorf("failed creating system deployer init transaction %w", err)
	}

	receipt, err := s.storage.GetFilteredInternalReceipt(context.Background(), tx.Hash(), nil, true)
	if err != nil {
		s.logger.Error("Load: Failed fetching receipt", "transactionHash", batch.Transactions[0].Hash().Hex(), "error", err)
		return fmt.Errorf("failed fetching receipt %w", err)
	}

	addresses, err := DeriveAddresses(receipt.ToReceipt())
	if err != nil {
		s.logger.Error("Load: Failed deriving addresses", "error", err, "receiptHash", receipt.TxHash.Hex())
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) initializeRequiredAddresses(addresses SystemContractAddresses) error {
	if addresses["TransactionsPostProcessor"] == nil {
		return fmt.Errorf("required contract address TransactionsPostProcessor is nil")
	}

	s.transactionsPostProcessorAddress = addresses["TransactionsPostProcessor"]
	s.systemAddresses = addresses

	return nil
}

func (s *systemContractCallbacks) Initialize(batch *core.Batch, receipt types.Receipt, msgBusManager SystemContractsInitializable) error {
	s.logger.Info("Initialize: Starting initialization of system contracts", "batchSeqNo", batch.SeqNo())
	if batch.SeqNo().Uint64() != common.L2SysContractGenesisSeqNo {
		s.logger.Error("Initialize: Batch is not genesis", "batchSeqNo", batch.SeqNo)
		return fmt.Errorf("batch is not genesis")
	}

	s.logger.Debug("Initialize: Deriving addresses from receipt", "transactionHash", receipt.TxHash.Hex())
	addresses, err := DeriveAddresses(&receipt)
	if err != nil {
		s.logger.Error("Initialize: Failed deriving addresses", "error", err, "receiptHash", receipt.TxHash.Hex())
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	if err := msgBusManager.Initialize(addresses); err != nil {
		s.logger.Error("Initialize: Failed deriving message bus address", "error", err)
		return fmt.Errorf("failed deriving message bus address %w", err)
	}

	s.logger.Info("Initialize: Initializing required addresses", "addresses", addresses)
	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) CreatePublicCallbackHandlerTransaction(ctx context.Context, l2State *state.StateDB) (*types.Transaction, error) {
	if s.PublicCallbackHandler() == nil {
		s.logger.Debug("CreatePublicCallbackHandlerTransaction: PublicCallbackHandler is nil, skipping transaction creation")
		return nil, nil
	}

	nonceForSyntheticTx := l2State.GetNonce(common.MaskedSender(*s.PublicCallbackHandler()))
	s.logger.Debug("CreatePublicCallbackHandlerTransaction: Retrieved nonce for synthetic transaction", "nonce", nonceForSyntheticTx)

	data, err := publicCallbacksABI.Pack("executeNextCallbacks")
	if err != nil {
		s.logger.Error("CreatePublicCallbackHandlerTransaction: Failed packing executeNextCallback data", "error", err)
		return nil, fmt.Errorf("failed packing executeNextCallback() %w", err)
	}

	tx := &types.LegacyTx{
		Nonce: nonceForSyntheticTx,
		Value: gethcommon.Big0,
		Gas:   params.MaxGasLimit,
		Data:  data,
		To:    s.PublicCallbackHandler(),
	}

	formedTx := types.NewTx(tx)
	s.logger.Info("CreatePublicCallbackHandlerTransaction: Successfully created transaction", "transactionHash", formedTx.Hash().Hex())
	return formedTx, nil
}

func (s *systemContractCallbacks) CreateOnBatchEndTransaction(_ context.Context, l2State *state.StateDB, results core.TxExecResults) (*types.Transaction, error) {
	if s.transactionsPostProcessorAddress == nil {
		s.logger.Debug("CreateOnBatchEndTransaction: TransactionsPostProcessorAddress is nil, skipping transaction creation")
		return nil, nil
	}

	if len(results) == 0 {
		s.logger.Debug("CreateOnBatchEndTransaction: Batch has no transactions, skipping transaction creation")
		return nil, ErrNoTransactions
	}

	nonceForSyntheticTx := l2State.GetNonce(common.MaskedSender(*s.transactionsPostProcessorAddress))
	s.logger.Debug("CreateOnBatchEndTransaction: Retrieved nonce for synthetic transaction", "nonce", nonceForSyntheticTx)

	// the data that is passed when the block ends
	synTxs := make([]TransactionPostProcessor.StructsTransaction, 0)
	for _, txExecResult := range results {
		tx := txExecResult.TxWithSender.Tx
		receipt := txExecResult.Receipt
		synTx := TransactionPostProcessor.StructsTransaction{
			Nonce:      big.NewInt(int64(txExecResult.TxWithSender.Tx.Nonce())),
			GasPrice:   tx.GasPrice(),
			GasLimit:   big.NewInt(int64(tx.Gas())),
			Value:      tx.Value(),
			Data:       tx.Data(),
			Successful: receipt.Status == types.ReceiptStatusSuccessful,
			GasUsed:    receipt.GasUsed,
		}
		if tx.To() != nil {
			synTx.To = *tx.To()
		} else {
			synTx.To = gethcommon.Address{} // Zero address - contract deployment
		}

		sender, err := core.GetExternalTxSigner(tx)
		if err != nil {
			s.logger.Error("CreateOnBatchEndTransaction: Failed to recover sender address", "error", err, "transactionHash", tx.Hash().Hex())
			return nil, fmt.Errorf("failed to recover sender address: %w", err)
		}
		synTx.From = sender

		synTxs = append(synTxs, synTx)
		s.logger.Debug("CreateOnBatchEndTransaction: Encoded transaction", log.TxKey, tx.Hash(), "sender", sender.Hex())
	}

	data, err := transactionPostProcessorABI.Pack("onBlock", synTxs)
	if err != nil {
		s.logger.Error("CreateOnBatchEndTransaction: Failed packing onBlock data", "error", err)
		return nil, fmt.Errorf("failed packing onBlock() %w", err)
	}

	tx := &types.LegacyTx{
		Nonce:    nonceForSyntheticTx,
		Value:    gethcommon.Big0,
		Gas:      params.MaxGasLimit,
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     data,
		To:       s.transactionsPostProcessorAddress,
	}

	formedTx := types.NewTx(tx)
	s.logger.Info("CreateOnBatchEndTransaction: Successfully created syntethic transaction", log.TxKey, formedTx.Hash())
	return formedTx, nil
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
