package enclave

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func decryptTransactions(txs nodecommon.EncryptedTransactions) Transactions {
	t := make([]L2Tx, 0)
	for _, tx := range txs {
		t = append(t, DecryptTx(tx))
	}
	return t
}

func DecryptTx(tx nodecommon.EncryptedTx) L2Tx {
	t := L2Tx{}
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		panic("no way")
	}

	return t
}

func EncryptTx(tx *L2Tx) nodecommon.EncryptedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic("no!")
	}
	return bytes
}

func encryptTransactions(transactions Transactions) nodecommon.EncryptedTransactions {
	result := make([]nodecommon.EncryptedTx, 0)
	for i := range transactions {
		result = append(result, EncryptTx(&transactions[i]))
	}
	return result
}

func DecryptRollup(rollup *nodecommon.Rollup) *Rollup {
	return &Rollup{
		Header:       rollup.Header,
		Transactions: decryptTransactions(rollup.Transactions),
	}
}
