package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SessionKeyExpirationService runs in the background and monitors users
type SessionKeyExpirationService struct {
	storage         storage.UserStorage
	logger          gethlog.Logger
	stopControl     *stopcontrol.StopControl
	ticker          *time.Ticker
	config          *wecommon.Config
	backendRPC      *BackendRPC
	activityTracker SessionKeyActivityTracker
	txSender        TxSender
}

// withSK opens an encrypted RPC connection authorized by the session key at `addr`
// and runs `fn`. Assumes user.SessionKeys[addr] exists.
func (s *SessionKeyExpirationService) withSK(
	ctx context.Context,
	user *wecommon.GWUser,
	addr common.Address,
	fn func(ctx context.Context, c *tenrpc.EncRPCClient) error,
) error {
	sk, ok := user.SessionKeys[addr]
	if !ok {
		return fmt.Errorf("session key not found for address %s", addr.Hex())
	}
	_, err := WithEncRPCConnection(ctx, s.backendRPC, sk.Account, func(c *tenrpc.EncRPCClient) (*struct{}, error) {
		return &struct{}{}, fn(ctx, c)
	})
	return err
}

// NewSessionKeyExpirationService creates a new session key expiration service
func NewSessionKeyExpirationService(storage storage.UserStorage, logger gethlog.Logger, stopControl *stopcontrol.StopControl, config *wecommon.Config, backendRPC *BackendRPC, activityTracker SessionKeyActivityTracker, txSender TxSender) *SessionKeyExpirationService {
	logger.Info("Creating session key expiration service", "expirationThreshold", config.SessionKeyExpirationThreshold.String())

	service := &SessionKeyExpirationService{
		storage:         storage,
		logger:          logger,
		stopControl:     stopControl,
		config:          config,
		backendRPC:      backendRPC,
		activityTracker: activityTracker,
		txSender:        txSender,
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
func (s *SessionKeyExpirationService) start() {
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
func (s *SessionKeyExpirationService) sessionKeyExpiration() {
	s.logger.Info("Session key expiration check started")

	cutoff := time.Now().Add(-s.config.SessionKeyExpirationThreshold)
	candidates := s.activityTracker.ListOlderThan(cutoff)

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
			s.activityTracker.Delete(c.Addr)
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

		_, err = s.txSender.SendAllMinusGasWithSK(context.Background(), user, c.Addr, *firstAccount.Address)
		if err != nil {
			s.logger.Error("Failed to recover funds from expired session key",
				"error", err,
				"userID", wecommon.HashForLogging(user.ID),
				"sessionKeyAddress", c.Addr.Hex())
			continue
		}

		// After successful external operation, delete from tracker
		_ = s.activityTracker.Delete(c.Addr)
	}
}
