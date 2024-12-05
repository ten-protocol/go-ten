package crosschain

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/system"
)

type (
	EVMExecutorResponse = map[common.TxHash]interface{}
	EVMExecutorFunc     = func(common.L2Transactions) map[common.TxHash]interface{}
	ObsCallEVMFunc      = func(core.Message) (*core.ExecutionResult, error)
)

type BlockMessageExtractor interface {
	// StoreCrossChainMessages - Verifies receipts belong to block and saves the relevant cross chain messages from the receipts
	StoreCrossChainMessages(ctx context.Context, block *types.Header, receipts common.L1Receipts) error

	StoreCrossChainValueTransfers(ctx context.Context, block *types.Header, receipts common.L1Receipts) error

	// GetBusAddress - Returns the L1 message bus address.
	GetBusAddress() *common.L1Address

	// Enabled - Returns true if there is a configured message bus, otherwise it is considered disabled
	Enabled() bool
}

type Manager interface {
	// IsSyntheticTransaction - Determines if a given L2 transaction is coming from the synthetic owner address.
	IsSyntheticTransaction(transaction *common.L2Tx) bool

	// GetOwner - Returns the address of the identity owning the message bus.
	GetOwner() common.L2Address

	// GetBusAddress - Returns the L2 address of the message bus contract.
	GetBusAddress() *common.L2Address

	// Initialize - Derives the address of the message bus contract.
	Initialize(systemAddresses system.SystemContractAddresses) error

	// GenerateMessageBusDeployTx - Returns a signed message bus deployment transaction.
	GenerateMessageBusDeployTx() (*common.L2Tx, error)

	// ExtractOutboundMessages - Finds relevant logs in the receipts and converts them to cross chain messages.
	ExtractOutboundMessages(ctx context.Context, receipts common.L2Receipts) (common.CrossChainMessages, error)

	ExtractOutboundTransfers(ctx context.Context, receipts common.L2Receipts) (common.ValueTransferEvents, error)

	CreateSyntheticTransactions(ctx context.Context, messages common.CrossChainMessages, stateDB *state.StateDB) common.L2Transactions

	ExecuteValueTransfers(ctx context.Context, transfers common.ValueTransferEvents, stateDB *state.StateDB)

	RetrieveInboundMessages(ctx context.Context, fromBlock *types.Header, toBlock *types.Header) (common.CrossChainMessages, common.ValueTransferEvents)

	system.SystemContractsInitializable
}
