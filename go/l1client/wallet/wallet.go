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
	SignTransaction(chainID int, tx *types.LegacyTx) (*types.Transaction, error)
}

type InMemoryWallet struct {
	pk *ecdsa.PrivateKey
}

func NewInMemoryWallet() Wallet {
	privateKey, _ := crypto.HexToECDSA("5dbbff1b5ff19f1ad6ea656433be35f6846e890b3f3ec6ef2b2e2137a8cab4ae")
	return &InMemoryWallet{
		pk: privateKey,
	}
}

func (m *InMemoryWallet) SignTransaction(chainID int, tx *types.LegacyTx) (*types.Transaction, error) {
	return types.SignNewTx(m.pk, types.NewEIP155Signer(big.NewInt(int64(chainID))), tx)
}

func (m *InMemoryWallet) Address() common.Address {
	publicKeyECDSA, ok := m.pk.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
