package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

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
		log.Error(fmt.Sprintf("could not decrypt encrypted L2 transaction. Cause: %s", err))
		panic(err)
	}

	return t
}

func EncryptTx(tx *nodecommon.L2Tx) nodecommon.EncryptedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Error(fmt.Sprintf("could not encrypt L2 transaction. Cause: %s", err))
		panic(err)
	}
	return bytes
}
