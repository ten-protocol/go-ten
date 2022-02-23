package enclave

import (
	"github.com/ethereum/go-ethereum/rlp"
	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
)

func decryptTransactions(txs common2.EncryptedTransactions) Transactions {
	t := make([]L2Tx, 0)
	for _, tx := range txs {
		t = append(t, DecryptTx(tx))
	}
	return t
}

func DecryptTx(tx common2.EncryptedTx) L2Tx {
	t := L2Tx{}
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		panic("no way")
	}

	return t
}

func EncryptTx(tx L2Tx) common2.EncryptedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic("no!")
	}
	return bytes
}

func encryptTransactions(transactions Transactions) common2.EncryptedTransactions {
	result := make([]common2.EncryptedTx, 0)
	for _, tx := range transactions {
		result = append(result, EncryptTx(tx))
	}
	return result
}

func DecryptRollup(rollup *common2.Rollup) *Rollup {
	return &Rollup{
		Header:       rollup.Header,
		Transactions: decryptTransactions(rollup.Transactions),
	}
}
