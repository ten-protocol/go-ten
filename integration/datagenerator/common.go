package datagenerator

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
)

func randomBytes(length int) []byte {
	byteArr := make([]byte, length)
	if _, err := rand.Read(byteArr); err != nil {
		panic(err)
	}

	return byteArr
}

func randomHash() common.Hash {
	return common.BytesToHash(randomBytes(32))
}

func randomAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

func randomUInt64() uint64 {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
	if err != nil {
		panic(err)
	}
	return val.Uint64()
}

func randomWithdrawal() nodecommon.Withdrawal {
	return nodecommon.Withdrawal{
		Amount:  randomUInt64(),
		Address: randomAddress(),
	}
}

func randomWithdrawals(length int) []nodecommon.Withdrawal {
	withdrawals := make([]nodecommon.Withdrawal, length)
	for i := 0; i < length; i++ {
		withdrawals[i] = randomWithdrawal()
	}
	return withdrawals
}

func randomEncryptedTransaction() nodecommon.EncryptedTx {
	return randomBytes(100)
}

func randomEncryptedTransactions(length int) nodecommon.EncryptedTransactions {
	encTransactions := make([]nodecommon.EncryptedTx, length)
	for i := 0; i < length; i++ {
		encTransactions[i] = randomEncryptedTransaction()
	}
	return encTransactions
}
