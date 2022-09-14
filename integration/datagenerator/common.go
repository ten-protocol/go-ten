package datagenerator

import (
	"crypto/rand"
	"github.com/ethereum/go-ethereum"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

func RandomBytes(length int) []byte {
	byteArr := make([]byte, length)
	if _, err := rand.Read(byteArr); err != nil {
		panic(err)
	}

	return byteArr
}

func randomHash() gethcommon.Hash {
	return gethcommon.BytesToHash(RandomBytes(32))
}

func RandomAddress() gethcommon.Address {
	return gethcommon.BytesToAddress(RandomBytes(20))
}

func randomUInt64() uint64 {
	val, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	return val.Uint64()
}

func randomWithdrawal() common.Withdrawal {
	return common.Withdrawal{
		Amount:    common.ValueInWei(big.NewInt(int64(randomUInt64()))),
		Recipient: RandomAddress(),
	}
}

func randomWithdrawals(length int) []common.Withdrawal {
	withdrawals := make([]common.Withdrawal, length)
	for i := 0; i < length; i++ {
		withdrawals[i] = randomWithdrawal()
	}
	return withdrawals
}

// CreateL2Tx Creates a dummy L2Tx for testing
func CreateL2Tx() *common.L2Tx {
	return types.NewTx(CreateL2TxData())
}

// CreateL2TxData Creates a dummy types.LegacyTx for testing
func CreateL2TxData() *types.LegacyTx {
	nonce, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	encodedTxData := make([]byte, 0)
	return &types.LegacyTx{
		Nonce: nonce.Uint64(), Value: big.NewInt(1), Gas: 1, GasPrice: big.NewInt(1), Data: encodedTxData,
	}
}

// CreateCallMsg Creates a dummy ethereum.CallMsg for testing
func CreateCallMsg() *ethereum.CallMsg {
	to := RandomAddress()
	return &ethereum.CallMsg{
		From:       RandomAddress(),
		To:         &to,
		Gas:        randomUInt64(),
		GasPrice:   big.NewInt(int64(randomUInt64())),
		GasFeeCap:  nil,
		GasTipCap:  nil,
		Value:      big.NewInt(int64(randomUInt64())),
		Data:       make([]byte, 0),
		AccessList: nil,
	}
}
