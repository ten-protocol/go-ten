package system

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/contracts/generated/TransactionsAnalyzer"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/wallet"
)

var (
	transactionsAnalyzerABI, _ = abi.JSON(strings.NewReader(TransactionsAnalyzer.TransactionsAnalyzerMetaData.ABI))
)

type SystemContractCallbacks interface {
	GetOwner() gethcommon.Address
	Initialize(batch *core.Batch, receipts types.Receipts) error
	Load() error
	CreateOnBatchEndTransaction(ctx context.Context, l2State *state.StateDB, batch *core.Batch, receipts common.L2Receipts) (*common.L2Tx, error)
	TransactionAnalyzerAddress() *gethcommon.Address
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

func (s *systemContractCallbacks) TransactionAnalyzerAddress() *gethcommon.Address {
	return s.transactionsAnalyzerAddress
}

func (s *systemContractCallbacks) GetOwner() gethcommon.Address {
	return s.ownerWallet.Address()
}

func (s *systemContractCallbacks) Load() error {
	if s.storage == nil {
		return fmt.Errorf("storage is not set")
	}

	batch, err := s.storage.FetchBatchBySeqNo(context.Background(), 1)
	if err != nil {
		return fmt.Errorf("failed fetching batch %w", err)
	}

	if len(batch.Transactions) < 2 {
		return fmt.Errorf("genesis batch does not have enough transactions")
	}

	receipt, err := s.storage.GetTransactionReceipt(context.Background(), batch.Transactions[1].Hash())
	if err != nil {
		return fmt.Errorf("failed fetching receipt %w", err)
	}

	addresses, err := DeriveAddresses(receipt)
	if err != nil {
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
	if len(receipts) < 2 {
		return fmt.Errorf("genesis batch does not have enough receipts")
	}

	addresses, err := DeriveAddresses(receipts[1])
	if err != nil {
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) CreateOnBatchEndTransaction(ctx context.Context, l2State *state.StateDB, batch *core.Batch, receipts common.L2Receipts) (*common.L2Tx, error) {
	if s.transactionsAnalyzerAddress == nil {
		return nil, nil
	}

	nonceForSyntheticTx := l2State.GetNonce(s.GetOwner())

	blockTransactions := TransactionsAnalyzer.TransactionsAnalyzerBlockTransactions{
		Transactions: make([][]byte, 0),
	}
	for _, tx := range batch.Transactions {

		encodedBytes, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return nil, fmt.Errorf("failed encoding transaction for onBlock %w", err)
		}

		blockTransactions.Transactions = append(blockTransactions.Transactions, encodedBytes)
	}

	data, err := transactionsAnalyzerABI.Pack("onBlock", blockTransactions)
	if err != nil {
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

	signedTx, err := s.ownerWallet.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed signing transaction %w", err)
	}

	return signedTx, nil
}
