package core

import (
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := datagenerator.CreateL2Tx()
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	tx1 := nodecommon.L2Tx{}
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
	rollup := datagenerator.RandomRollup()
	_, read, err := rlp.EncodeToReader(&rollup)
	if err != nil {
		panic(err)
	}
	r1 := nodecommon.Rollup{}

	err = rlp.Decode(read, &r1)
	if err != nil {
		panic(err)
	}
	if r1.Hash() != rollup.Hash() {
		t.Errorf("rollup deserialized incorrectly\n")
	}
}
