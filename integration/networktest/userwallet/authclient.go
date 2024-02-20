package userwallet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration"
)

const (
	_maxReceiptWaitTime  = 30 * time.Second
	_receiptPollInterval = 1 * time.Second // todo (@matt) this should be configured using network timings provided by env
)

// AuthClientUser is a test user that uses the auth client to talk to directly to a node
// Note: AuthClientUser is **not** thread-safe for a single wallet (creates nonce conflicts etc.)
type AuthClientUser struct {
	privateKey     *ecdsa.PrivateKey
	publicKey      *ecdsa.PublicKey
	accountAddress gethcommon.Address
	chainID        *big.Int
	rpcEndpoint    string

	// state managed by the wallet
	nonce uint64

	client *obsclient.AuthObsClient // lazily initialised and authenticated on first usage
	logger gethlog.Logger
}

// Option modifies a AuthClientUser. See below for options, in the form `WithXxx(xxx)` that can be chained into constructor
type Option func(wallet *AuthClientUser)

func NewUserWallet(pk *ecdsa.PrivateKey, rpcEndpoint string, logger gethlog.Logger, opts ...Option) *AuthClientUser {
	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		// this shouldn't happen
		logger.Crit("error casting public key to ECDSA")
	}
	wal := &AuthClientUser{
		privateKey:     pk,
		publicKey:      publicKeyECDSA,
		accountAddress: crypto.PubkeyToAddress(*publicKeyECDSA),
		chainID:        big.NewInt(integration.TenChainID), // default, overridable using `WithChainID(...) opt`
		rpcEndpoint:    rpcEndpoint,
		logger:         logger,
	}
	// apply any optional config to the wallet
	for _, opt := range opts {
		opt(wal)
	}
	return wal
}

func (s *AuthClientUser) SendFunds(ctx context.Context, addr gethcommon.Address, value *big.Int) (*gethcommon.Hash, error) {
	err := s.EnsureClientSetup(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare client to send funds - %w", err)
	}

	txData := &types.LegacyTx{
		Nonce: s.nonce,
		Value: value,
		To:    &addr,
	}
	tx := s.client.EstimateGasAndGasPrice(txData) //nolint: contextcheck

	txHash, err := s.SendTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("unable to send transaction - %w", err)
	}

	return txHash, nil
}

func (s *AuthClientUser) SendTransaction(ctx context.Context, tx types.TxData) (*gethcommon.Hash, error) {
	signedTx, err := s.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("unable to sign transaction - %w", err)
	}
	// fmt.Printf("waiting for receipt hash %s\n", signedTx.Hash())
	err = s.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("unable to send transaction - %w", err)
	}

	txHash := signedTx.Hash()
	// transaction has been sent, we increment the nonce
	s.nonce++
	return &txHash, nil
}

func (s *AuthClientUser) AwaitReceipt(ctx context.Context, txHash *gethcommon.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	err = retry.Do(func() error {
		receipt, err = s.client.TransactionReceipt(ctx, *txHash)
		if !errors.Is(err, rpc.ErrNilResponse) {
			// nil response means not found. Any other error is unexpected, so we stop polling and fail immediately
			return retry.FailFast(err)
		}
		return err
	}, retry.NewTimeoutStrategy(_maxReceiptWaitTime, _receiptPollInterval))
	return receipt, err
}

func (s *AuthClientUser) Address() gethcommon.Address {
	return s.accountAddress
}

func (s *AuthClientUser) SignTransaction(tx types.TxData) (*types.Transaction, error) {
	return s.SignTransactionForChainID(tx, s.chainID)
}

func (s *AuthClientUser) SignTransactionForChainID(tx types.TxData, chainID *big.Int) (*types.Transaction, error) {
	return types.SignNewTx(s.privateKey, types.NewLondonSigner(chainID), tx)
}

func (s *AuthClientUser) GetNonce() uint64 {
	return s.nonce
}

func (s *AuthClientUser) PrivateKey() *ecdsa.PrivateKey {
	return s.privateKey
}

//
// These methods allow the user to comply with the wallet.Wallet interface
//

func (s *AuthClientUser) ChainID() *big.Int {
	return s.chainID
}

func (s *AuthClientUser) SetNonce(_ uint64) {
	panic("AuthClientUser is designed to manage its own nonce - this method exists to support legacy interface methods")
}

func (s *AuthClientUser) GetNonceAndIncrement() uint64 {
	panic("AuthClientUser is designed to manage its own nonce - this method exists to support legacy interface methods")
}

// EnsureClientSetup creates an authenticated RPC client (with a viewing key generated, signed and registered) when first called
// Also fetches current nonce value.
func (s *AuthClientUser) EnsureClientSetup(ctx context.Context) error {
	if s.client != nil {
		// client already setup
		return nil
	}
	authClient, err := obsclient.DialWithAuth(s.rpcEndpoint, s, s.logger)
	if err != nil {
		return err
	}
	s.client = authClient

	// fetch current nonce for account
	nonce, err := authClient.NonceAt(ctx, big.NewInt(-1))
	if err != nil {
		return fmt.Errorf("unable to fetch client nonce - %w", err)
	}
	s.nonce = nonce

	return nil
}

func (s *AuthClientUser) NativeBalance(ctx context.Context) (*big.Int, error) {
	err := s.EnsureClientSetup(ctx)
	if err != nil {
		return nil, err
	}
	return s.client.BalanceAt(ctx, nil)
}

// Init forces VK setup: currently the faucet http server requires a viewing key for a wallet to even *receive* funds :(
func (s *AuthClientUser) Init(ctx context.Context) (*AuthClientUser, error) {
	return s, s.EnsureClientSetup(ctx)
}

// UserWalletOptions can be passed into the constructor to override default values
// e.g. NewUserWallet(pk, rpcAddr, logger, WithChainId(123))
// NewUserWallet(pk, rpcAddr, logger, WithChainId(123), WithRPCTimeout(20*time.Second)), )

func WithChainID(chainID *big.Int) Option {
	return func(wallet *AuthClientUser) {
		wallet.chainID = chainID
	}
}
