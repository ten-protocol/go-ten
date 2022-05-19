package core

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/core/types"
)

// Creates a dummy L2Tx for testing
func CreateL2Tx() *nodecommon.L2Tx {
	return types.NewTx(CreateL2TxData())
}

// Creates a dummy types.LegacyTx for testing
func CreateL2TxData() *types.LegacyTx {
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	encodedTxData := make([]byte, 0)
	return &types.LegacyTx{
		Nonce: nonce.Uint64(), Value: big.NewInt(1), Gas: 1, GasPrice: big.NewInt(1), Data: encodedTxData,
	}
}
