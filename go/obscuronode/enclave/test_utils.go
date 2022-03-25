package enclave

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Creates a dummy L2Tx for testing
func createL2Tx() *nodecommon.L2Tx {
	return types.NewTx(createL2TxData())
}

// Creates a dummy types.LegacyTx for testing
func createL2TxData() *types.LegacyTx {
	txData := L2TxData{TransferTx, common.Address{}, common.Address{}, 100}
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	encodedTxData, _ := rlp.EncodeToBytes(txData)
	return &types.LegacyTx{
		Nonce: nonce.Uint64(), Value: big.NewInt(1), Gas: 1, GasPrice: big.NewInt(1), Data: encodedTxData,
	}
}
