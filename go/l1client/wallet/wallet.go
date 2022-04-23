package wallet

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet interface {
	Address() common.Address
	SignTransaction(chainID int, tx types.TxData) (*types.Transaction, error)
}

type inMemoryWallet struct {
	prvKey     *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	pubKeyAddr common.Address
}

func NewInMemoryWallet(pk string) Wallet {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(err)
	}
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return &inMemoryWallet{
		prvKey:     privateKey,
		pubKey:     publicKeyECDSA,
		pubKeyAddr: crypto.PubkeyToAddress(*publicKeyECDSA),
	}
}

// SignTransaction returns a signed transaction
func (m *inMemoryWallet) SignTransaction(chainID int, tx types.TxData) (*types.Transaction, error) {
	return types.SignNewTx(m.prvKey, types.NewEIP155Signer(big.NewInt(int64(chainID))), tx)
}

// Address returns the current wallet address
func (m *inMemoryWallet) Address() common.Address {
	return m.pubKeyAddr
}
