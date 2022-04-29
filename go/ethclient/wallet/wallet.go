package wallet

import (
	"crypto/ecdsa"
	"math/big"
	"sync"

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
	SetNonce(nonce uint64)
	// Nonce returns the current nonce
	Nonce() uint64
	// IncrementNonce increments the nonce by one
	IncrementNonce()
}

type inMemoryWallet struct {
	prvKey     *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	pubKeyAddr common.Address
	nonce      uint64
	nlock      sync.RWMutex
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

func (m *inMemoryWallet) Nonce() uint64 {
	m.nlock.RLock()
	defer m.nlock.RUnlock()
	return m.nonce
}

func (m *inMemoryWallet) IncrementNonce() {
	m.nlock.Lock()
	defer m.nlock.Unlock()
	m.nonce++
}

func (m *inMemoryWallet) SetNonce(nonce uint64) {
	m.nlock.Lock()
	defer m.nlock.Unlock()
	m.nonce = nonce
}
