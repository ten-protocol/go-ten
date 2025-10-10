package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

		// Transfer funds to user's primary account using TxSender (sends all minus gas)
		// Find the first account registered with the user - we will send funds to this account
		var firstAccount *wecommon.GWAccount
		for _, account := range user.Accounts {
			firstAccount = account
			break
		}
		if firstAccount == nil || firstAccount.Address == nil {
			s.logger.Error("No primary account found for user", "userID", wecommon.HashForLogging(user.ID))
			continue
		}

		_, err = s.services.TxSender.SendAllMinusGasWithSK(context.Background(), user, c.Addr, *firstAccount.Address)
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
