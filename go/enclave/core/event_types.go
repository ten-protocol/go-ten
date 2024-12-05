package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventVisibilityConfig - configuration per event by the dApp developer(DD)
// There are 4 cases:
// 1. DD doesn't configure anything. - ContractVisibilityConfig.AutoConfig=true
// 2. DD configures and  specifies the contract as transparent - ContractVisibilityConfig.Transparent=true
// 3. DD configures and specify the contract as non-transparent, but doesn't configure the event - Contract: false/false , EventVisibilityConfig.AutoConfig=true
// DD configures the contract as non-transparent, and also configures the topics for the event
type EventVisibilityConfig struct {
	AutoConfig                                  bool  // true for events that have no explicit configuration
	Public                                      bool  // everyone can see and query for this event
	Topic1CanView, Topic2CanView, Topic3CanView *bool // If the event is not public, and this is true, it means that the address from topicI is an EOA that can view this event
	SenderCanView                               *bool // if true, the tx signer will see this event. Default false
}

// ContractVisibilityConfig represents the configuration as defined by the dApp developer in the smart contract
type ContractVisibilityConfig struct {
	AutoConfig   bool                                       // true for contracts that have no explicit configuration
	Transparent  *bool                                      // users can configure contracts to be fully transparent. All events will be public, and it will expose the internal storage.
	EventConfigs map[gethcommon.Hash]*EventVisibilityConfig // map from the event log signature (topics[0]) to the settings
}

type TxExecResult struct {
	Receipt          *types.Receipt
	CreatedContracts map[gethcommon.Address]*ContractVisibilityConfig
	TxWithSender     *TxWithSender
	Err              error
}

type TxWithSender struct {
	Tx          *types.Transaction
	Sender      *gethcommon.Address
	IsSynthetic bool
}

type TransactionsWithSender []*TxWithSender

func (stxs *TransactionsWithSender) Add(tx *common.L2PricedTransaction) error {
	sender, err := GetTxSigner(tx)
	if err != nil {
		return err
	}
	*stxs = append(*stxs, &TxWithSender{Tx: tx.Tx, Sender: &sender})
	return nil
}

// InternalReceipt - Equivalent to the geth types.Receipt, but without weird quirks
type InternalReceipt struct {
	PostState         []byte
	Status            uint64
	CumulativeGasUsed uint64
	EffectiveGasPrice *uint64
	CreatedContract   *gethcommon.Address
	TxHash            gethcommon.Hash
	BlockHash         gethcommon.Hash
	BlockNumber       *big.Int
	TransactionIndex  uint
	From              gethcommon.Address
	To                *gethcommon.Address
	TxType            uint8
	Logs              []*types.Log
}

// MarshalToJson marshals a transaction receipt into a JSON object.
// taken from geth
func (receipt *InternalReceipt) MarshalToJson() map[string]interface{} {
	var effGasPrice *hexutil.Big
	if receipt.EffectiveGasPrice != nil {
		effGasPrice = (*hexutil.Big)(big.NewInt(int64(*receipt.EffectiveGasPrice)))
	}

	fields := map[string]interface{}{
		"blockHash":         receipt.BlockHash,
		"blockNumber":       hexutil.Uint64(receipt.BlockNumber.Uint64()),
		"transactionHash":   receipt.TxHash,
		"transactionIndex":  hexutil.Uint64(receipt.TransactionIndex),
		"from":              receipt.From,
		"to":                receipt.To,
		"gasUsed":           hexutil.Uint64(receipt.CumulativeGasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   receipt.CreatedContract,
		"logs":              receipt.Logs,
		"logsBloom":         types.Bloom{},
		"type":              hexutil.Uint(receipt.TxType),
		"effectiveGasPrice": effGasPrice,
	}

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = []*types.Log{}
	}

	return fields
}

func (receipt *InternalReceipt) ToReceipt() *types.Receipt {
	var effGasPrice *big.Int
	if receipt.EffectiveGasPrice != nil {
		effGasPrice = big.NewInt(int64(*receipt.EffectiveGasPrice))
	}

	var cc gethcommon.Address
	if receipt.CreatedContract != nil {
		cc = *receipt.CreatedContract
	}
	return &types.Receipt{
		Type:              receipt.TxType,
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		Bloom:             types.Bloom{},
		Logs:              receipt.Logs,
		TxHash:            receipt.TxHash,
		ContractAddress:   cc,
		GasUsed:           receipt.CumulativeGasUsed,
		EffectiveGasPrice: effGasPrice,
		BlobGasUsed:       0,
		BlobGasPrice:      nil,
		BlockHash:         receipt.BlockHash,
		BlockNumber:       receipt.BlockNumber,
		TransactionIndex:  receipt.TransactionIndex,
	}
}

type TxExecResults []*TxExecResult

func (txResults *TxExecResults) Add(other ...*TxExecResult) {
	*txResults = append(*txResults, other...)
}

func (txResults *TxExecResults) MarkSynthetic(isSynthetic bool) {
	for _, txResult := range *txResults {
		txResult.TxWithSender.IsSynthetic = isSynthetic
	}
}

func (txResults *TxExecResults) GetSynthetic() *TxExecResults {
	syntheticTxs := make(TxExecResults, 0)
	for _, txResult := range *txResults {
		if txResult.TxWithSender.IsSynthetic {
			syntheticTxs = append(syntheticTxs, txResult)
		}
	}
	return &syntheticTxs
}

func (txResults *TxExecResults) GetReal() *TxExecResults {
	realTxs := make(TxExecResults, 0)
	for _, txResult := range *txResults {
		if !txResult.TxWithSender.IsSynthetic {
			realTxs = append(realTxs, txResult)
		}
	}
	return &realTxs
}

func (txResults *TxExecResults) ToTransactionsWithSenders() TransactionsWithSender {
	transactionsWithSenders := make(TransactionsWithSender, len(*txResults))
	for i, txResult := range *txResults {
		transactionsWithSenders[i] = txResult.TxWithSender
	}
	return transactionsWithSenders
}

func (txResults *TxExecResults) Transactions() []*common.L2Tx {
	txs := make([]*common.L2Tx, len(*txResults))
	for i, txResult := range *txResults {
		txs[i] = txResult.TxWithSender.Tx
	}
	return txs
}

func (txResults *TxExecResults) Receipts() types.Receipts {
	receipts := make(types.Receipts, len(*txResults))
	for i, txResult := range *txResults {
		receipts[i] = txResult.Receipt
	}
	return receipts
}
