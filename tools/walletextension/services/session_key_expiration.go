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
			s.logger.Debug("Checking session key",
				"userID", wecommon.HashForLogging(user.ID),
				"sessionKeyAddress", sessionKeyAddr.Hex(),
				"createdAt", sessionKey.CreatedAt.Format(time.RFC3339),
				"age", age.String(),
				"expirationThreshold", expirationThreshold.String(),
				"isExpired", age > expirationThreshold)
			
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
				} else {
					// If expired session key has non-zero balance, send funds back to first account
					s.logger.Info("Expired session key with non-zero balance - attempting fund recovery",
						"userID", wecommon.HashForLogging(user.ID),
						"sessionKeyAddress", sessionKeyAddr.Hex(),
						"balance", balanceStr)

					err := s.recoverFundsFromExpiredSessionKey(user, sessionKeyAddr, balance)
					if err != nil {
						s.logger.Error("Failed to recover funds from expired session key",
							"error", err,
							"userID", wecommon.HashForLogging(user.ID),
							"sessionKeyAddress", sessionKeyAddr.Hex(),
							"balance", balanceStr)
					} else {
						s.logger.Info("Successfully initiated fund recovery from expired session key",
							"userID", wecommon.HashForLogging(user.ID),
							"sessionKeyAddress", sessionKeyAddr.Hex(),
							"balance", balanceStr)
					}
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

// recoverFundsFromExpiredSessionKey sends funds from an expired session key to the user's first account
func (s *SessionKeyExpirationService) recoverFundsFromExpiredSessionKey(user *wecommon.GWUser, sessionKeyAddr common.Address, balance *hexutil.Big) error {
	ctx := context.Background()

	// Find the first account registered with the user
	var firstAccount *wecommon.GWAccount
	for _, account := range user.Accounts {
		firstAccount = account
		break // Get the first account
	}

	if firstAccount == nil {
		return fmt.Errorf("no accounts found for user %s", wecommon.HashForLogging(user.ID))
	}

	s.logger.Info("Recovering funds from expired session key",
		"userID", wecommon.HashForLogging(user.ID),
		"sessionKeyAddress", sessionKeyAddr.Hex(),
		"recipientAddress", firstAccount.Address.Hex(),
		"amount", balance.String())

	// Use the existing transaction API to send funds
	// This leverages the existing SendTransaction infrastructure that handles session keys properly
	txHash, err := s.sendFundsUsingTransactionAPI(ctx, user, sessionKeyAddr, *firstAccount.Address, balance)
	if err != nil {
		return fmt.Errorf("failed to send funds transaction: %w", err)
	}

	s.logger.Info("Fund recovery transaction sent",
		"userID", wecommon.HashForLogging(user.ID),
		"sessionKeyAddress", sessionKeyAddr.Hex(),
		"recipientAddress", firstAccount.Address.Hex(),
		"amount", balance.String(),
		"transactionHash", txHash.Hex())

	return nil
}

// sendFundsUsingTransactionAPI uses the existing transaction API to send funds from a session key
func (s *SessionKeyExpirationService) sendFundsUsingTransactionAPI(ctx context.Context, user *wecommon.GWUser, fromAddr, toAddr common.Address, amount *hexutil.Big) (common.Hash, error) {
	s.logger.Info("Attempting to send funds using transaction API",
		"userID", wecommon.HashForLogging(user.ID),
		"fromAddr", fromAddr.Hex(),
		"toAddr", toAddr.Hex(),
		"amount", amount.String())

	// Create transaction args that will be handled by the existing SendTransaction API
	// The transaction API will automatically detect this is a session key transaction and handle it properly
	txArgs := map[string]interface{}{
		"from":  fromAddr.Hex(),
		"to":    toAddr.Hex(),
		"value": amount.String(),
		"gas":   "0x5208", // 21000 gas for simple transfer
	}

	s.logger.Debug("Transaction args created", "txArgs", txArgs)

	// Use the existing RPC infrastructure to call the transaction API
	// This will go through the same path as if a dApp sent the transaction
	result, err := WithEncRPCConnection(ctx, s.services.BackendRPC, user.SessionKeys[fromAddr].Account, func(rpcClient *tenrpc.EncRPCClient) (*common.Hash, error) {
		var txHash common.Hash
		s.logger.Debug("Calling eth_sendTransaction", "txArgs", txArgs)
		err := rpcClient.CallContext(ctx, &txHash, "eth_sendTransaction", txArgs)
		if err != nil {
			s.logger.Error("eth_sendTransaction failed", "error", err)
		} else {
			s.logger.Info("eth_sendTransaction succeeded", "txHash", txHash.Hex())
		}
		return &txHash, err
	})

	if err != nil {
		s.logger.Error("WithEncRPCConnection failed", "error", err)
		return common.Hash{}, err
	}

	s.logger.Info("Fund recovery transaction sent successfully", "txHash", result.Hex())
	return *result, nil
}
