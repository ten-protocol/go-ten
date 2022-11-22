package crosschain

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

type OnChainEVMExecutorResponse = map[common.TxHash]interface{}
type OnChainEVMExecutorFunc = func(common.L2Transactions) map[common.TxHash]interface{}
type OffChainEVMCallFunc = func(types.Message) (*core.ExecutionResult, error)

type MainNetMessageExtractor interface {
	//(todo: Verifies receipts match block) & extracts the messages logged inside the receipts
	ProcessCrossChainMessages(block *common.L1Block, receipts common.L1Receipts) error

	GetBusAddress() *common.L1Address

	Enabled() bool
}

type ObscuroCrossChainManager interface {
	//Address of the identity owning the message bus.
	GetOwner() common.L2Address

	//Address of the messageBusContract
	GetBusAddress() *common.L2Address

	//Generate the key pair that will be used to transact with the L2 message bus.
	DeriveOwner(seed []byte) (*common.L2Address, error) //TODO: Implement with cryptography epic.

	//Will create a message bus deployment transaction OR give an error on why it can't.
	GenerateMessageBusDeployTx() (*common.L2Tx, error)

	//
	ExtractLocalMessages(receipts common.L2Receipts) (common.CrossChainMessages, error)

	//Executes synthetic transactions
	SubmitRemoteMessagesLocally(fromBlock *common.L1Block, toBlock *common.L1Block, rollupState *state.StateDB, processTxCall OnChainEVMExecutorFunc, processOffChainMessage OffChainEVMCallFunc) error
}
