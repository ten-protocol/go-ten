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

// SessionKeyExpirationService runs in the background and monitors users
type SessionKeyExpirationService struct {
	storage     storage.UserStorage
	logger      gethlog.Logger
	stopControl *stopcontrol.StopControl
	ticker      *time.Ticker
	config      *wecommon.Config
	services    *Services // Reference to main services for RPC access
}

// NewSessionKeyExpirationService creates a new session key expiration service
func NewSessionKeyExpirationService(storage storage.UserStorage, logger gethlog.Logger, stopControl *stopcontrol.StopControl, config *wecommon.Config, services *Services) *SessionKeyExpirationService {
	service := &SessionKeyExpirationService{
		storage:     storage,
		logger:      logger,
		stopControl: stopControl,
		config:      config,
		services:    services,
	}

	// Start the service
	go service.start()

	return service
}

// start begins the periodic user monitoring
func (s *SessionKeyExpirationService) start() {
	// TODO: Make this interval configurable
	s.ticker = time.NewTicker(10 * time.Second)

	go func() {
		defer s.ticker.Stop()

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

	s.logger.Info("Session key expiration service started")
}

// sessionKeyExpiration runs the monitoring logic
func (s *SessionKeyExpirationService) sessionKeyExpiration() {
	users, err := s.storage.GetAllUsers()
	if err != nil {
		s.logger.Error("Failed to get all users", "error", err)
		return
	}

	s.logger.Info("Session key expiration service running - checking users", "count", len(users))

	// Use configurable expiration threshold
	expirationThreshold := s.config.SessionKeyExpirationThreshold
	now := time.Now()
	expiredKeysCount := 0

	for _, user := range users {
		s.logger.Info("User found",
			"userID", wecommon.HashForLogging(user.ID),
			"accountsCount", len(user.Accounts),
			"sessionKeysCount", len(user.SessionKeys))

		// Check each session key for expiration
		for sessionKeyAddr, sessionKey := range user.SessionKeys {
			age := now.Sub(sessionKey.CreatedAt)
			if age > expirationThreshold {
				expiredKeysCount++

				// Check balance for expired session key
				balance, err := s.getSessionKeyBalance(user, sessionKeyAddr)
				balanceStr := "unknown"
				hasZeroBalance := false
				if err != nil {
					s.logger.Error("Failed to get balance for expired session key",
						"error", err,
						"sessionKeyAddress", sessionKeyAddr.Hex())
				} else {
					balanceStr = balance.String()
					// Check if balance is zero
					hasZeroBalance = balance.ToInt().Cmp(big.NewInt(0)) == 0
				}

				s.logger.Warn("Expired session key found",
					"userID", wecommon.HashForLogging(user.ID),
					"sessionKeyAddress", sessionKeyAddr.Hex(),
					"createdAt", sessionKey.CreatedAt.Format(time.RFC3339),
					"age", age.String(),
					"expirationThreshold", expirationThreshold.String(),
					"balance", balanceStr)

				// If expired session key has zero balance, it can be deleted
				if hasZeroBalance {
					s.logger.Info("Expired session key with zero balance can be deleted",
						"userID", wecommon.HashForLogging(user.ID),
						"sessionKeyAddress", sessionKeyAddr.Hex(),
						"balance", balanceStr)

					// TODO: Uncomment the following code to actually delete the session key
					// err := s.storage.RemoveSessionKey(user.ID, &sessionKeyAddr)
					// if err != nil {
					// 	s.logger.Error("Failed to delete expired session key with zero balance",
					// 		"error", err,
					// 		"userID", wecommon.HashForLogging(user.ID),
					// 		"sessionKeyAddress", sessionKeyAddr.Hex())
					// } else {
					// 	s.logger.Info("Successfully deleted expired session key with zero balance",
					// 		"userID", wecommon.HashForLogging(user.ID),
					// 		"sessionKeyAddress", sessionKeyAddr.Hex())
					// }
				}
			}
		}
	}

	if expiredKeysCount > 0 {
		s.logger.Info("Session key expiration check completed",
			"totalExpiredKeys", expiredKeysCount,
			"expirationThreshold", expirationThreshold.String())
	} else {
		s.logger.Info("Session key expiration check completed - no expired keys found",
			"expirationThreshold", expirationThreshold.String())
	}
}

// getSessionKeyBalance retrieves the balance for a session key account
func (s *SessionKeyExpirationService) getSessionKeyBalance(user *wecommon.GWUser, sessionKeyAddr common.Address) (*hexutil.Big, error) {
	ctx := context.Background()

	// Use the latest block for balance checking
	latest := rpc.LatestBlockNumber
	blockNrOrHash := rpc.BlockNumberOrHash{
		BlockNumber: &latest,
	}

	// Create a temporary user with only the session key account for authorization
	tempUser := &wecommon.GWUser{
		ID:          user.ID,
		Accounts:    make(map[common.Address]*wecommon.GWAccount),
		UserKey:     user.UserKey,
		SessionKeys: make(map[common.Address]*wecommon.GWSessionKey),
	}

	// Add the specific session key to the temp user
	if sessionKey, exists := user.SessionKeys[sessionKeyAddr]; exists {
		tempUser.SessionKeys[sessionKeyAddr] = sessionKey
	} else {
		return nil, fmt.Errorf("session key not found for address %s", sessionKeyAddr.Hex())
	}

	// Use the existing WithEncRPCConnection infrastructure to get balance
	// This is similar to how the blockchain API works but without the import cycle
	balance, err := WithEncRPCConnection(ctx, s.services.BackendRPC, tempUser.SessionKeys[sessionKeyAddr].Account, func(rpcClient *tenrpc.EncRPCClient) (*hexutil.Big, error) {
		var result hexutil.Big
		err := rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCGetBalance, sessionKeyAddr, blockNrOrHash)
		return &result, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get balance via RPC: %w", err)
	}

	return balance, nil
}
