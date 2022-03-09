package enclave

import (
	"crypto/rand"
	"math"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/rlp"
	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := createL2Tx()
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	tx1 := L2Tx{}
	err2 := rlp.DecodeBytes(bytes, &tx1)
	if err2 != nil {
		panic(err2)
	}
	if tx1.Hash() != tx.Hash() {
		t.Errorf("tx deserialized incorrectly\n")
	}
}

func TestSerialiseRollup(t *testing.T) {
	tx := createL2Tx()
	height := atomic.Value{}
	height.Store(1)
	rollup := common2.Rollup{
		Header:       GenesisRollup.Header,
		Height:       height,
		Transactions: encryptTransactions(Transactions{*tx}),
	}
	_, read, err := rlp.EncodeToReader(&rollup)
	if err != nil {
		panic(err)
	}
	r1 := common2.Rollup{}

	err = rlp.Decode(read, &r1)
	if err != nil {
		panic(err)
	}
	if r1.Hash() != rollup.Hash() {
		t.Errorf("rollup deserialized incorrectly\n")
	}
}

// Creates a dummy L2Tx for testing
func createL2Tx() *L2Tx {
	txData := L2TxData{TransferTx, common.Address{}, common.Address{}, 100}
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	encodedTxData, _ := rlp.EncodeToBytes(txData)
	return types.NewTx(&types.LegacyTx{
		Nonce: nonce.Uint64(), Value: big.NewInt(1), Gas: 1, GasPrice: big.NewInt(1), Data: encodedTxData,
	})
}
