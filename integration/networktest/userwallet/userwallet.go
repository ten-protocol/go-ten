package userwallet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/networktest"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
)

const (
	_maxReceiptWaitTime  = 30 * time.Second
	_receiptPollInterval = 1 * time.Second // todo (@matt) this should be configured using network timings provided by env
)

// UserWallet implements wallet.Wallet so it can be used with the original Wallet code.
// But it aims to provide a wider range of functionality, akin to the software and hardware wallets that users interact with.
// Note: UserWallet is **not** thread-safe for a single wallet (creates nonce conflicts etc.)
type UserWallet struct {
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

// Option modifies a UserWallet. See below for options, in the form `WithXxx(xxx)` that can be chained into constructor
type Option func(wallet *UserWallet)

// GenerateRandomWallet will generate a random wallet with a UserWallet wrapper, connecting to a random validator node
// Note: will use testlog.Logger() as the logger
func GenerateRandomWallet(network networktest.NetworkConnector) *UserWallet {
	wallet := datagenerator.RandomWallet(network.ChainID())
	_, err := obsclient.DialWithAuth(network.SequencerRPCAddress(), wallet, testlog.Logger())
	if err != nil {
		panic(err)
	}

	rndValidatorIdx := int(datagenerator.RandomUInt64()) % network.NumValidators()
	return NewUserWallet(wallet.PrivateKey(), network.ValidatorRPCAddress(rndValidatorIdx), testlog.Logger())
}

func NewUserWallet(pk *ecdsa.PrivateKey, rpcEndpoint string, logger gethlog.Logger, opts ...Option) *UserWallet {
	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		// this shouldn't happen
		logger.Crit("error casting public key to ECDSA")
	}
	wal := &UserWallet{
		privateKey:     pk,
		publicKey:      publicKeyECDSA,
		accountAddress: crypto.PubkeyToAddress(*publicKeyECDSA),
		chainID:        big.NewInt(integration.ObscuroChainID), // default, overridable using `WithChainID(...) opt`
		rpcEndpoint:    rpcEndpoint,
		logger:         logger,
	}
	// apply any optional config to the wallet
	for _, opt := range opts {
		opt(wal)
	}
	return wal
}

func (s *UserWallet) ChainID() *big.Int {
	return big.NewInt(integration.ObscuroChainID)
}

func (s *UserWallet) SendFunds(ctx context.Context, addr gethcommon.Address, value *big.Int) (*gethcommon.Hash, error) {
	err := s.EnsureClientSetup(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare client to send funds - %w", err)
	}

	tx := &types.LegacyTx{
		Nonce:    s.nonce,
		Value:    value,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		To:       &addr,
	}

	txHash, err := s.SendTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("unable to send transaction - %w", err)
	}

	return txHash, nil
}

func (s *UserWallet) SendTransaction(ctx context.Context, tx *types.LegacyTx) (*gethcommon.Hash, error) {
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

func (s *UserWallet) AwaitReceipt(ctx context.Context, txHash *gethcommon.Hash) (*types.Receipt, error) {
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

func (s *UserWallet) Address() gethcommon.Address {
	return s.accountAddress
}

func (s *UserWallet) SignTransaction(tx types.TxData) (*types.Transaction, error) {
	return types.SignNewTx(s.privateKey, types.NewLondonSigner(s.chainID), tx)
}

func (s *UserWallet) GetNonce() uint64 {
	return s.nonce
}

func (s *UserWallet) PrivateKey() *ecdsa.PrivateKey {
	return s.privateKey
}

func (s *UserWallet) SetNonce(_ uint64) {
	panic("UserWallet is designed to manage its own nonce - this method exists to support legacy interface methods")
}

func (s *UserWallet) GetNonceAndIncrement() uint64 {
	panic("UserWallet is designed to manage its own nonce - this method exists to support legacy interface methods")
}

// EnsureClientSetup creates an authenticated RPC client (with a viewing key generated, signed and registered) when first called
// Also fetches current nonce value.
func (s *UserWallet) EnsureClientSetup(ctx context.Context) error {
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
	nonce, err := authClient.NonceAt(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to fetch client nonce - %w", err)
	}
	s.nonce = nonce

	return nil
}

// ResetClient creates an authenticated RPC client (with a viewing key generated, signed and registered)
// Also fetches current nonce value. It closes previous client if it exists.
func (s *UserWallet) ResetClient(ctx context.Context) error {
	if s.client != nil {
		// client already setup, close it before re-authenticating
		s.client.Close()
	}
	authClient, err := obsclient.DialWithAuth(s.rpcEndpoint, s, s.logger)
	if err != nil {
		return err
	}
	s.client = authClient

	// fetch current nonce for account
	nonce, err := authClient.NonceAt(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to fetch client nonce - %w", err)
	}
	s.nonce = nonce

	return nil
}

func (s *UserWallet) NativeBalance(ctx context.Context) (*big.Int, error) {
	err := s.EnsureClientSetup(ctx)
	if err != nil {
		return nil, err
	}
	return s.client.BalanceAt(ctx, nil)
}

// Init forces VK setup: currently the faucet http server requires a viewing key for a wallet to even *receive* funds :(
func (s *UserWallet) Init(ctx context.Context) (*UserWallet, error) {
	return s, s.EnsureClientSetup(ctx)
}

// UserWalletOptions can be passed into the constructor to override default values
// e.g. NewUserWallet(pk, rpcAddr, logger, WithChainId(123))
// NewUserWallet(pk, rpcAddr, logger, WithChainId(123), WithRPCTimeout(20*time.Second)), )

func WithChainID(chainID *big.Int) Option {
	return func(wallet *UserWallet) {
		wallet.chainID = chainID
	}
}
