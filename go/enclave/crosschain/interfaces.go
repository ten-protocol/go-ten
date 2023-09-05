package crosschain

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/obscuronet/go-obscuro/go/common"
)

type (
	EVMExecutorResponse = map[common.TxHash]interface{}
	EVMExecutorFunc     = func(common.L2Transactions) map[common.TxHash]interface{}
	ObsCallEVMFunc      = func(core.Message) (*core.ExecutionResult, error)
)

type BlockMessageExtractor interface {
	// StoreCrossChainMessages - Verifies receipts belong to block and saves the relevant cross chain messages from the receipts
	StoreCrossChainMessages(block *common.L1Block, receipts common.L1Receipts) error

	StoreCrossChainValueTransfers(block *common.L1Block, receipts common.L1Receipts) error

	// GetBusAddress - Returns the L1 message bus address.
	GetBusAddress() *common.L1Address

	// Enabled - Returns true if there is a configured message bus, otherwise it is considered disabled
	Enabled() bool
}

type Manager interface {
	// IsSyntheticTransaction - Determines if a given L2 transaction is coming from the synthetic owner address.
	IsSyntheticTransaction(transaction common.L2Tx) bool

	// GetOwner - Returns the address of the identity owning the message bus.
	GetOwner() common.L2Address

	// GetBusAddress - Returns the L2 address of the message bus contract.
	GetBusAddress() *common.L2Address

	// DeriveOwner - Generates the key pair that will be used to transact with the L2 message bus.
	// todo (#1549) - implement with cryptography epic.
	DeriveOwner(seed []byte) (*common.L2Address, error)

	// GenerateMessageBusDeployTx - Returns a signed message bus deployment transaction.
	GenerateMessageBusDeployTx() (*common.L2Tx, error)

	// ExtractOutboundMessages - Finds relevant logs in the receipts and converts them to cross chain messages.
	ExtractOutboundMessages(receipts common.L2Receipts) (common.CrossChainMessages, error)

	ExtractOutboundTransfers(receipts common.L2Receipts) (common.ValueTransferEvents, error)

	CreateSyntheticTransactions(messages common.CrossChainMessages, rollupState *state.StateDB) common.L2Transactions

	ExecuteValueTransfers(transfers common.ValueTransferEvents, rollupState *state.StateDB)

	RetrieveInboundMessages(fromBlock *common.L1Block, toBlock *common.L1Block, rollupState *state.StateDB) (common.CrossChainMessages, common.ValueTransferEvents)
}
