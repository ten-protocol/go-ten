package wallet

import (
	"crypto/ecdsa"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet interface {
	// Address returns the pubkey Address of the wallet
	Address() common.Address
	// SignTransaction returns a signed transaction
	SignTransaction(chainID int, tx types.TxData) (*types.Transaction, error)

	// SetNonce overrides the current nonce
	// The Nonce is expected to be the next nonce to use in a transaction, not the current account Nonce
	SetNonce(nonce uint64)
	// GetNonceAndIncrement atomically increments the nonce by one and returns the previous value
	GetNonceAndIncrement() uint64
}

type inMemoryWallet struct {
	prvKey     *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	pubKeyAddr common.Address
	nonce      uint64
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

func (m *inMemoryWallet) GetNonceAndIncrement() uint64 {
	return atomic.AddUint64(&m.nonce, 1) - 1
}

func (m *inMemoryWallet) SetNonce(nonce uint64) {
	m.nonce = nonce
}
