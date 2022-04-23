package wallet

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet interface {
	Address() common.Address
	SignTransaction(chainID int, tx types.TxData) (*types.Transaction, error)
}

type InMemoryWallet struct {
	pk *ecdsa.PrivateKey
}

func NewInMemoryWallet(pk string) Wallet {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(err)
	}
	return &InMemoryWallet{
		pk: privateKey,
	}
}

func (m *InMemoryWallet) SignTransaction(chainID int, tx types.TxData) (*types.Transaction, error) {
	return types.SignNewTx(m.pk, types.NewEIP155Signer(big.NewInt(int64(chainID))), tx)
}

func (m *InMemoryWallet) Address() common.Address {
	publicKeyECDSA, ok := m.pk.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
