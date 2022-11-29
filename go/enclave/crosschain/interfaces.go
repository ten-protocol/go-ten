package crosschain

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

type (
	OnChainEVMExecutorResponse = map[common.TxHash]interface{}
	OnChainEVMExecutorFunc     = func(common.L2Transactions) map[common.TxHash]interface{}
	OffChainEVMCallFunc        = func(types.Message) (*core.ExecutionResult, error)
)

type BlockMessageExtractor interface {
	// ProcessCrossChainMessages - Verifies receipts belong to block and saves the relevant cross chain messages from the receipts
	ProcessCrossChainMessages(block *common.L1Block, receipts common.L1Receipts) error

	// GetBusAddress - Returns the L1 message bus address.
	GetBusAddress() *common.L1Address

	// Enabled - Returns true if there is a configured message bus, otherwise it is considered disabled
	Enabled() bool
}

type ObscuroCrossChainManager interface {
	// GetOwner - Returns the address of the identity owning the message bus.
	GetOwner() common.L2Address

	// GetBusAddress - Returns the L2 address of the message bus contract.
	GetBusAddress() *common.L2Address

	// DeriveOwner - Generates the key pair that will be used to transact with the L2 message bus.
	// TODO: Implement with cryptography epic.
	DeriveOwner(seed []byte) (*common.L2Address, error)

	// GenerateMessageBusDeployTx - Returns a signed message bus deployment transaction.
	GenerateMessageBusDeployTx() (*common.L2Tx, error)

	//ExtractLocalMessages - Finds relevant logs in the receipts and converts them to cross chain messages.
	ExtractLocalMessages(receipts common.L2Receipts) (common.CrossChainMessages, error)

	// SubmitRemoteMessagesLocally - Submits messages saved between the from and to blocks on chain using the provided function bindings.
	SubmitRemoteMessagesLocally(fromBlock *common.L1Block, toBlock *common.L1Block, rollupState *state.StateDB, processTxCall OnChainEVMExecutorFunc, processOffChainMessage OffChainEVMCallFunc) error
}
