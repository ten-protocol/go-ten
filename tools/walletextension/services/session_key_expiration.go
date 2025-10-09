package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	tencommonrpc "github.com/ten-protocol/go-ten/go/common/rpc"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// dustThresholdWei is the minimum balance considered worth recovering from an expired
// session key. Set in wei for precision.
// // 1_000_000_000_000 wei = 1e12 wei = 1,000 gwei = 0.000001 ETH (~$0.004 at $4,000/ETH)
// Adjust the USD intuition based on current ETH price.
var dustThresholdWei = big.NewInt(1_000_000_000_000)

// sessionKeyExpirationService runs in the background and monitors users
type sessionKeyExpirationService struct {
	storage     storage.UserStorage
	logger      gethlog.Logger
	stopControl *stopcontrol.StopControl
	ticker      *time.Ticker
	config      *wecommon.Config
	services    *Services // Reference to main services for RPC access
}

// withSK opens an encrypted RPC connection authorized by the session key at `addr`
// and runs `fn`. Assumes user.SessionKeys[addr] exists.
func (s *sessionKeyExpirationService) withSK(
	ctx context.Context,
	user *wecommon.GWUser,
	addr common.Address,
	fn func(ctx context.Context, c *tenrpc.EncRPCClient) error,
) error {
	sk, ok := user.SessionKeys[addr]
	if !ok {
		return fmt.Errorf("session key not found for address %s", addr.Hex())
	}
	_, err := WithEncRPCConnection(ctx, s.services.BackendRPC, sk.Account, func(c *tenrpc.EncRPCClient) (*struct{}, error) {
		return &struct{}{}, fn(ctx, c)
	})
	return err
}

// NewSessionKeyExpirationService creates a new session key expiration service
func NewSessionKeyExpirationService(storage storage.UserStorage, logger gethlog.Logger, stopControl *stopcontrol.StopControl, config *wecommon.Config, services *Services) *sessionKeyExpirationService {
	logger.Info("Creating session key expiration service", "expirationThreshold", config.SessionKeyExpirationThreshold.String())

	service := &sessionKeyExpirationService{
		storage:     storage,
		logger:      logger,
		stopControl: stopControl,
		config:      config,
		services:    services,
	}

	// Start the service
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Session key expiration service creation panicked", "error", r)
			}
		}()
		service.start()
	}()
	return service
}

// start begins the periodic user monitoring
func (s *sessionKeyExpirationService) start() {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Session key expiration service panicked", "error", r)
		}
	}()

	// Configure interval from config
	interval := s.config.SessionKeyExpirationInterval
	s.ticker = time.NewTicker(interval)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Session key expiration goroutine panicked", "error", r)
			}
			s.ticker.Stop()
		}()

		for {
			select {
			case <-s.ticker.C:
				s.sessionKeyExpiration()
			case <-s.stopControl.Done():
				s.logger.Info("Session key expiration service stopped")
				return
			}
		}
	}()
	s.logger.Info("Session key expiration service started", "expirationThreshold", s.config.SessionKeyExpirationThreshold.String(), "interval", interval.String())
}

// sessionKeyExpiration runs the monitoring logic
func (s *sessionKeyExpirationService) sessionKeyExpiration() {
	s.logger.Info("Session key expiration check started")

	cutoff := time.Now().Add(-s.config.SessionKeyExpirationThreshold)
	candidates := s.services.ActivityTracker.ListOlderThan(cutoff)

	if len(candidates) == 0 {
		return
	}

	for _, c := range candidates {
		// Load the user for this session key
		user, err := s.storage.GetUser(c.UserID)
		if err != nil || user == nil {
			s.logger.Error("Failed to load user for session key candidate", "error", err)
			continue
		}

		// Ensure this session key still belongs to the user
		if _, ok := user.SessionKeys[c.Addr]; !ok {
			// The session key may have been deleted; remove from tracker
			s.services.ActivityTracker.Delete(c.Addr)
			continue
		}

		// Check balance for expired session key
		balance, err := s.getSessionKeyBalance(user, c.Addr)
		if err != nil {
			s.logger.Error("Failed to get balance for expired session key", "error", err, "sessionKeyAddress", c.Addr.Hex())
			continue
		}

		if balance.ToInt().Cmp(dustThresholdWei) <= 0 {
			continue
		}

		// Transfer funds to user's primary account
		err = s.transferExpiredSessionKeyFundsToPrimaryAccount(user, c.Addr, balance)
		if err != nil {
			s.logger.Error("Failed to recover funds from expired session key",
				"error", err,
				"userID", wecommon.HashForLogging(user.ID),
				"sessionKeyAddress", c.Addr.Hex(),
				"balance", balance)
			continue
		}

		// After successful external operation, delete from tracker
		_ = s.services.ActivityTracker.Delete(c.Addr)
	}
}

// transferExpiredSessionKeyFundsToPrimaryAccount sends funds from an expired session key to the user's first account
func (s *sessionKeyExpirationService) transferExpiredSessionKeyFundsToPrimaryAccount(user *wecommon.GWUser, sessionKeyAddr common.Address, balance *hexutil.Big) error {
	ctx := context.Background()

	// Find the first account registered with the user - we will send funds to this account
	var firstAccount *wecommon.GWAccount
	for _, account := range user.Accounts {
		firstAccount = account
		break // Get the first account
	}

	if firstAccount == nil {
		return fmt.Errorf("no accounts found for user %s", wecommon.HashForLogging(user.ID))
	}

	gasPrice, err := s.getGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	gasLimit, err := s.estimateGas(ctx, user, sessionKeyAddr, *firstAccount.Address, balance)
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Calculate gas cost: gasPrice * gasLimit
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	balanceInt := balance.ToInt()

	// Calculate amount to send: balance - gasCost
	amountToSend := new(big.Int).Sub(balanceInt, gasCost)

	legacyTx := &types.LegacyTx{
		To:       firstAccount.Address,
		Value:    amountToSend,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	}

	tx := types.NewTx(legacyTx)
	if tx == nil {
		return fmt.Errorf("failed to create transaction")
	}

	// Sign the transaction with the session key
	signedTx, err := s.services.SKManager.SignTx(ctx, user, sessionKeyAddr, tx)
	if err != nil {
		s.logger.Error("Failed to sign transaction with session key", "error", err, "sessionKeyAddress", sessionKeyAddr.Hex())
		return fmt.Errorf("failed to sign transaction with session key: %w", err)
	}
	blob, err := signedTx.MarshalBinary()
	if err != nil {
		s.logger.Error("Failed to marshal signed transaction", "error", err)
		return fmt.Errorf("failed to marshal signed transaction: %w", err)
	}

	_, err = s.sendRawTransaction(ctx, blob, user, sessionKeyAddr)
	if err != nil {
		s.logger.Error("Failed to send transaction", "error", err, "sessionKeyAddress", sessionKeyAddr.Hex())
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	// we don't wait for the receipt here if there was an error we will retry this transaction in the next inverval

	return nil
}

// getGasPrice retrieves the current gas price
func (s *sessionKeyExpirationService) getGasPrice(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	_, err := WithPlainRPCConnection(ctx, s.services.BackendRPC, func(client *rpc.Client) (*hexutil.Big, error) {
		err := client.CallContext(ctx, &result, tenrpc.GasPrice)
		return &result, err
	})
	if err != nil {
		s.logger.Error("Failed to get gas price", "error", err)
		return nil, fmt.Errorf("failed to get gas price via RPC: %w", err)
	}
	gasPrice := result.ToInt()
	return gasPrice, nil
}

// estimateGas estimates gas for a transfer
func (s *sessionKeyExpirationService) estimateGas(ctx context.Context, user *wecommon.GWUser, from, to common.Address, value *hexutil.Big) (uint64, error) {
	// Session key presence validated in withSK

	var result hexutil.Uint64
	err := s.withSK(ctx, user, from, func(ctx context.Context, rpcClient *tenrpc.EncRPCClient) error {
		params := map[string]interface{}{
			"from":  from.Hex(),
			"to":    to.Hex(),
			"value": value.String(),
		}
		return rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCEstimateGas, params)
	})
	if err != nil {
		s.logger.Error("Failed to estimate gas", "error", err, "from", from.Hex(), "to", to.Hex(), "value", value.String())
		return 0, fmt.Errorf("failed to estimate gas via RPC: %w", err)
	}
	return uint64(result), nil
}

// sendRawTransaction sends a raw transaction using the same approach as SendRawTx
func (s *sessionKeyExpirationService) sendRawTransaction(ctx context.Context, input hexutil.Bytes, user *wecommon.GWUser, sessionKeyAddr common.Address) (common.Hash, error) {
	// Session key presence validated in withSK

	var result common.Hash
	err := s.withSK(ctx, user, sessionKeyAddr, func(ctx context.Context, rpcClient *tenrpc.EncRPCClient) error {
		return rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCSendRawTransaction, input)
	})
	if err != nil {
		return common.Hash{}, err
	}
	return result, nil
}

// getSessionKeyBalance retrieves the balance for a session key account
func (s *sessionKeyExpirationService) getSessionKeyBalance(user *wecommon.GWUser, sessionKeyAddr common.Address) (*hexutil.Big, error) {
	ctx := context.Background()

	// Use the pending block for balance checking so pending txs are reflected
	pending := rpc.PendingBlockNumber
	blockNrOrHash := rpc.BlockNumberOrHash{
		BlockNumber: &pending,
	}

	// Use the EncRPC to get balance via session key auth
	var balance *hexutil.Big
	err := s.withSK(ctx, user, sessionKeyAddr, func(ctx context.Context, rpcClient *tenrpc.EncRPCClient) error {
		var result hexutil.Big
		if callErr := rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCGetBalance, sessionKeyAddr, blockNrOrHash); callErr != nil {
			return callErr
		}
		balance = &result
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get balance via RPC: %w", err)
	}

	return balance, nil
}
