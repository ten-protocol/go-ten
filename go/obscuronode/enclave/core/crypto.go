package core

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

func DecodeTransactions(txs nodecommon.EncodedTransactions) L2Txs {
	t := make([]nodecommon.L2Tx, 0)
	for _, tx := range txs {
		t = append(t, DecodeTx(tx))
	}
	return t
}

func DecodeTx(tx nodecommon.EncodedTx) nodecommon.L2Tx {
	t := nodecommon.L2Tx{}
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		log.Panic("could not decrypt encrypted L2 transaction. Cause: %s", err)
	}

	return t
}

func EncodeTx(tx *nodecommon.L2Tx) nodecommon.EncodedTx {
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Panic("could not encrypt L2 transaction. Cause: %s", err)
	}
	return bytes
}
