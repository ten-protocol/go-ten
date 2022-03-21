package enclave

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
)

type (
	L2Txs []nodecommon.L2Tx
	// L2TxType indicates the type of L2 transaction - either a transfer or a withdrawal for now
	L2TxType uint8
)

const (
	TransferTx L2TxType = iota
	WithdrawalTx
)

// L2TxData is the Obscuro transaction data that will be stored encoded in the types.Transaction data field.
type L2TxData struct {
	Type   L2TxType
	From   common.Address
	To     common.Address
	Amount uint64
}

// TxData returns the decoded L2 data stored in the transaction's data field.
func TxData(tx *nodecommon.L2Tx) L2TxData {
	data := L2TxData{}

	err := rlp.DecodeBytes(tx.Data(), &data)
	if err != nil {
		// TODO - Surface this error properly.
		panic(err)
	}

	return data
}
