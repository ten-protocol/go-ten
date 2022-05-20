package wallet

import (
	"crypto/ecdsa"
	"math/big"
	"sync/atomic"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet interface {
	// Address returns the pubkey Address of the wallet
	Address() common.Address
	// SignTransaction returns a signed transaction
	SignTransaction(tx types.TxData) (*types.Transaction, error)

	// SetNonce overrides the current nonce
	// The Nonce is expected to be the next nonce to use in a transaction, not the current account Nonce
	SetNonce(nonce uint64)
	// GetNonceAndIncrement atomically increments the nonce by one and returns the previous value
	GetNonceAndIncrement() uint64

	// PK is a temp method that should be removed when the l2 wallet (wallet_mock.Wallet) is updated
	// TODO remove this
	PK() *ecdsa.PrivateKey
}

type inMemoryWallet struct {
	prvKey     *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	pubKeyAddr common.Address
	nonce      uint64
	chainID    *big.Int
}

func NewInMemoryWallet(chainID *big.Int, pk string) Wallet {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Panic("could not recover private key from hex. Cause: %s", err)
	}
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Panic("error casting public key to ECDSA")
	}

	return &inMemoryWallet{
		chainID:    chainID,
		prvKey:     privateKey,
		pubKey:     publicKeyECDSA,
		pubKeyAddr: crypto.PubkeyToAddress(*publicKeyECDSA),
	}
}

// SignTransaction returns a signed transaction
func (m *inMemoryWallet) SignTransaction(tx types.TxData) (*types.Transaction, error) {
	return types.SignNewTx(m.prvKey, types.NewEIP155Signer(m.chainID), tx)
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

func (m *inMemoryWallet) PK() *ecdsa.PrivateKey {
	return m.prvKey
}
