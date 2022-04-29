package core

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func DecryptTransactions(txs nodecommon.EncryptedTransactions) L2Txs {
	t := make([]nodecommon.L2Tx, 0)
	for _, tx := range txs {
		t = append(t, DecryptTx(tx))
	}
	return t
}

func DecryptTx(tx nodecommon.EncryptedTx) nodecommon.L2Tx {
	t := nodecommon.L2Tx{}
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		panic(err)
	}

	return t
}

func EncryptTx(tx *nodecommon.L2Tx) nodecommon.EncryptedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic("no!")
	}
	return bytes
}

//func DecryptRollup(rollup *nodecommon.Rollup) *Rollup {
//	return &Rollup{
//		Header:       rollup.Header,
//		Transactions: DecryptTransactions(rollup.Transactions),
//	}
//}

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte
