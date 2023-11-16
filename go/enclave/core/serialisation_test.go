package core

import (
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := datagenerator.CreateL2Tx()
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	tx1 := common.L2Tx{}
	err2 := rlp.DecodeBytes(bytes, &tx1)
	if err2 != nil {
		panic(err2)
	}
	if tx1.Hash() != tx.Hash() {
		t.Errorf("tx deserialized incorrectly\n")
	}
}

func TestSerialiseRollup(t *testing.T) {
	height := atomic.Value{}
	height.Store(1)
	rollup := datagenerator.RandomRollup(nil)
	_, read, err := rlp.EncodeToReader(&rollup)
	if err != nil {
		panic(err)
	}
	r1 := common.ExtRollup{}

	err = rlp.Decode(read, &r1)
	if err != nil {
		panic(err)
	}
	if r1.Hash() != rollup.Hash() {
		t.Errorf("rollup deserialized incorrectly\n")
	}
}

func TestSerialiseBatch(t *testing.T) {
	height := atomic.Value{}
	height.Store(1)
	batch := datagenerator.RandomBatch(nil)
	_, read, err := rlp.EncodeToReader(&batch)
	if err != nil {
		panic(err)
	}
	r1 := common.ExtBatch{}

	err = rlp.Decode(read, &r1)
	if err != nil {
		panic(err)
	}
	if r1.Hash() != batch.Hash() {
		t.Errorf("batch deserialized incorrectly\n")
	}
}
