package userwallet //nolint:typecheck

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/wallet"
)

const (
	_maxReceiptWaitTime  = 30 * time.Second
	_receiptPollInterval = 1 * time.Second // todo (@matt) this should be configured using network timings provided by env
)

// AuthClientUser is a test user that uses the auth client to talk to directly to a node
// Note: AuthClientUser is **not** thread-safe for a single wallet (creates nonce conflicts etc.)
type AuthClientUser struct {
	wal         wallet.Wallet
	rpcEndpoint string

	client *obsclient.AuthObsClient // lazily initialised and authenticated on first usage
	logger gethlog.Logger
}

func NewUserWallet(wal wallet.Wallet, rpcEndpoint string, logger gethlog.Logger) *AuthClientUser {
	return &AuthClientUser{
		wal:         wal,
		rpcEndpoint: rpcEndpoint,
		logger:      logger,
	}
}

func (s *AuthClientUser) SendFunds(ctx context.Context, addr gethcommon.Address, value *big.Int) (*gethcommon.Hash, error) {
	err := s.EnsureClientSetup(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare client to send funds - %w", err)
	}

	txData := &types.LegacyTx{
		Nonce: s.wal.GetNonce(),
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
	signedTx, err := s.wal.SignTransaction(tx)
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
	s.wal.GetNonceAndIncrement()
	return &txHash, nil
}

func (s *AuthClientUser) AwaitReceipt(ctx context.Context, txHash *gethcommon.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	err = retry.Do(func() error {
		receipt, err = s.client.TransactionReceipt(ctx, *txHash)
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			return retry.FailFast(err)
		}
		return err
	}, retry.NewTimeoutStrategy(_maxReceiptWaitTime, _receiptPollInterval))
	return receipt, err
}

// EnsureClientSetup creates an authenticated RPC client (with a viewing key generated, signed and registered) when first called
// Also fetches current nonce value.
func (s *AuthClientUser) EnsureClientSetup(ctx context.Context) error {
	if s.client != nil {
		// client already setup
		return nil
	}
	authClient, err := obsclient.DialWithAuth(s.rpcEndpoint, s.wal, s.logger)
	if err != nil {
		return err
	}
	s.client = authClient

	// fetch current nonce for account
	nonce, err := authClient.NonceAt(ctx, big.NewInt(-1))
	if err != nil {
		return fmt.Errorf("unable to fetch client nonce - %w", err)
	}
	s.wal.SetNonce(nonce)

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

func (s *AuthClientUser) Wallet() wallet.Wallet {
	return s.wal
}

func (s *AuthClientUser) GetPersonalTransactions(ctx context.Context, pagination common.QueryPagination) (types.Receipts, error) {
	address := s.wal.Address()
	return s.client.GetReceiptsByAddress(ctx, &address, pagination)
}
