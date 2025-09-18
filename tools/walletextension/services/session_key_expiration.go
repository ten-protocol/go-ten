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
	fmt.Printf("🔧 Creating session key expiration service with threshold: %s\n", config.SessionKeyExpirationThreshold.String())
	logger.Info("Creating session key expiration service", "expirationThreshold", config.SessionKeyExpirationThreshold.String())

	service := &SessionKeyExpirationService{
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
				fmt.Printf("❌ Session key expiration service creation panicked: %v\n", r)
				logger.Error("Session key expiration service creation panicked", "error", r)
			}
		}()
		service.start()
	}()

	fmt.Printf("✅ Session key expiration service created successfully\n")
	logger.Info("Session key expiration service created successfully")
	return service
}

// start begins the periodic user monitoring
func (s *SessionKeyExpirationService) start() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ Session key expiration service panicked: %v\n", r)
			s.logger.Error("Session key expiration service panicked", "error", r)
		}
	}()

	fmt.Printf("🚀 Session key expiration service starting with threshold: %s\n", s.config.SessionKeyExpirationThreshold.String())
	s.logger.Info("Session key expiration service starting", "expirationThreshold", s.config.SessionKeyExpirationThreshold.String())

	// TODO: Make this interval configurable
	s.ticker = time.NewTicker(10 * time.Second)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("❌ Session key expiration goroutine panicked: %v\n", r)
				s.logger.Error("Session key expiration goroutine panicked", "error", r)
			}
			s.ticker.Stop()
		}()

		for {
			select {
			case <-s.ticker.C:
				s.sessionKeyExpiration()
			case <-s.stopControl.Done():
				fmt.Printf("🛑 Session key expiration service stopped\n")
				s.logger.Info("Session key expiration service stopped")
				return
			}
		}
	}()

	fmt.Printf("✅ Session key expiration service started successfully\n")
	s.logger.Info("Session key expiration service started", "expirationThreshold", s.config.SessionKeyExpirationThreshold.String())
}

// sessionKeyExpiration runs the monitoring logic
func (s *SessionKeyExpirationService) sessionKeyExpiration() {
	fmt.Printf("🔍 Session key expiration check started\n")
	s.logger.Info("Session key expiration check started")

	users, err := s.storage.GetAllUsers()
	if err != nil {
		fmt.Printf("❌ Failed to get all users: %v\n", err)
		s.logger.Error("Failed to get all users", "error", err)
		return
	}

	// Use configurable expiration threshold
	expirationThreshold := s.config.SessionKeyExpirationThreshold

	fmt.Printf("👥 Session key expiration service running - checking %d users with threshold: %s\n", len(users), expirationThreshold.String())
	s.logger.Info("Session key expiration service running - checking users", "count", len(users), "expirationThreshold", expirationThreshold.String())
	now := time.Now()
	expiredKeysCount := 0

	for _, user := range users {
		fmt.Printf("👤 User found - ID: %s, Accounts: %d, SessionKeys: %d\n",
			wecommon.HashForLogging(user.ID), len(user.Accounts), len(user.SessionKeys))
		s.logger.Info("User found",
			"userID", wecommon.HashForLogging(user.ID),
			"accountsCount", len(user.Accounts),
			"sessionKeysCount", len(user.SessionKeys))

		// Check each session key for expiration
		for sessionKeyAddr, sessionKey := range user.SessionKeys {
			age := now.Sub(sessionKey.CreatedAt)
			isExpired := age > expirationThreshold
			fmt.Printf("🔑 Checking session key: %s, Created: %s, Age: %s, Expired: %v\n",
				sessionKeyAddr.Hex(), sessionKey.CreatedAt.Format(time.RFC3339), age.String(), isExpired)
			s.logger.Info("Checking session key",
				"userID", wecommon.HashForLogging(user.ID),
				"sessionKeyAddress", sessionKeyAddr.Hex(),
				"createdAt", sessionKey.CreatedAt.Format(time.RFC3339),
				"age", age.String(),
				"expirationThreshold", expirationThreshold.String(),
				"isExpired", isExpired)

			if age > expirationThreshold {
				expiredKeysCount++
				fmt.Printf("⏰ EXPIRED SESSION KEY FOUND: %s (age: %s)\n", sessionKeyAddr.Hex(), age.String())

				// Check balance for expired session key
				balance, err := s.getSessionKeyBalance(user, sessionKeyAddr)
				balanceStr := "unknown"
				hasZeroBalance := false
				if err != nil {
					fmt.Printf("❌ Failed to get balance for expired session key %s: %v\n", sessionKeyAddr.Hex(), err)
					s.logger.Error("Failed to get balance for expired session key",
						"error", err,
						"sessionKeyAddress", sessionKeyAddr.Hex())
				} else {
					balanceStr = balance.String()
					// Check if balance is zero
					hasZeroBalance = balance.ToInt().Cmp(big.NewInt(0)) == 0
					fmt.Printf("💰 Expired session key balance: %s (zero: %v)\n", balanceStr, hasZeroBalance)
				}

				fmt.Printf("⚠️  Expired session key found - User: %s, Key: %s, Balance: %s\n",
					wecommon.HashForLogging(user.ID), sessionKeyAddr.Hex(), balanceStr)
				s.logger.Warn("Expired session key found",
					"userID", wecommon.HashForLogging(user.ID),
					"sessionKeyAddress", sessionKeyAddr.Hex(),
					"createdAt", sessionKey.CreatedAt.Format(time.RFC3339),
					"age", age.String(),
					"expirationThreshold", expirationThreshold.String(),
					"balance", balanceStr)

				// If expired session key has zero balance, it can be deleted
				if hasZeroBalance {
					fmt.Printf("🗑️  Expired session key with zero balance can be deleted: %s\n", sessionKeyAddr.Hex())
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
					fmt.Printf("💸 Expired session key with non-zero balance - attempting fund recovery: %s (balance: %s)\n",
						sessionKeyAddr.Hex(), balanceStr)
					s.logger.Info("Expired session key with non-zero balance - attempting fund recovery",
						"userID", wecommon.HashForLogging(user.ID),
						"sessionKeyAddress", sessionKeyAddr.Hex(),
						"balance", balanceStr)

					err := s.recoverFundsFromExpiredSessionKey(user, sessionKeyAddr, balance)
					if err != nil {
						fmt.Printf("❌ Failed to recover funds from expired session key %s: %v\n", sessionKeyAddr.Hex(), err)
						s.logger.Error("Failed to recover funds from expired session key",
							"error", err,
							"userID", wecommon.HashForLogging(user.ID),
							"sessionKeyAddress", sessionKeyAddr.Hex(),
							"balance", balanceStr)
					} else {
						fmt.Printf("✅ Successfully initiated fund recovery from expired session key %s\n", sessionKeyAddr.Hex())
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

	fmt.Printf("🔄 Starting fund recovery from expired session key %s\n", sessionKeyAddr.Hex())

	// Find the first account registered with the user
	var firstAccount *wecommon.GWAccount
	for _, account := range user.Accounts {
		firstAccount = account
		break // Get the first account
	}

	if firstAccount == nil {
		fmt.Printf("❌ No accounts found for user %s\n", wecommon.HashForLogging(user.ID))
		return fmt.Errorf("no accounts found for user %s", wecommon.HashForLogging(user.ID))
	}

	fmt.Printf("💸 Recovering funds: %s from %s to %s\n",
		balance.String(), sessionKeyAddr.Hex(), firstAccount.Address.Hex())
	s.logger.Info("Recovering funds from expired session key",
		"userID", wecommon.HashForLogging(user.ID),
		"sessionKeyAddress", sessionKeyAddr.Hex(),
		"recipientAddress", firstAccount.Address.Hex(),
		"amount", balance.String())

	// Send funds to the first account that was authenticated with the gateway
	// using the same approach as the existing transaction API

	// Get gas price and estimate gas for the transfer
	fmt.Printf("⛽ Getting gas price...\n")
	gasPrice, err := s.getGasPrice(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get gas price: %v\n", err)
		return fmt.Errorf("failed to get gas price: %w", err)
	}
	fmt.Printf("⛽ Gas price: %s\n", gasPrice.String())

	// Estimate gas for the transfer (using the full balance initially)
	fmt.Printf("⛽ Estimating gas for transfer...\n")
	gasLimit, err := s.estimateGas(ctx, sessionKeyAddr, *firstAccount.Address, balance)
	if err != nil {
		fmt.Printf("❌ Failed to estimate gas: %v\n", err)
		return fmt.Errorf("failed to estimate gas: %w", err)
	}
	fmt.Printf("⛽ Gas limit: %d\n", gasLimit)

	// Calculate gas cost
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	fmt.Printf("⛽ Gas cost: %s\n", gasCost.String())

	// Calculate the amount to send (balance minus gas cost)
	amountToSend := new(big.Int).Sub(balance.ToInt(), gasCost)
	fmt.Printf("💰 Amount to send: %s (balance: %s - gas: %s)\n", amountToSend.String(), balance.String(), gasCost.String())

	// Check if we have enough balance to cover gas + transfer
	if amountToSend.Cmp(big.NewInt(0)) <= 0 {
		fmt.Printf("❌ Insufficient balance to cover gas costs: balance=%s, gasCost=%s\n", balance.String(), gasCost.String())
		s.logger.Warn("Insufficient balance to cover gas costs",
			"sessionKeyAddress", sessionKeyAddr.Hex(),
			"balance", balance.String(),
			"gasCost", gasCost.String(),
			"gasPrice", gasPrice.String(),
			"gasLimit", gasLimit)
		return fmt.Errorf("insufficient balance to cover gas costs: balance=%s, gasCost=%s", balance.String(), gasCost.String())
	}

	// Get nonce for the session key
	fmt.Printf("🔢 Getting nonce for session key...\n")
	nonce, err := s.getNonce(ctx, sessionKeyAddr)
	if err != nil {
		fmt.Printf("❌ Failed to get nonce: %v\n", err)
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	fmt.Printf("🔢 Nonce: %d\n", nonce)

	fmt.Printf("📋 Transaction details - From: %s, To: %s, Amount: %s, Gas: %s, Nonce: %d\n",
		sessionKeyAddr.Hex(), firstAccount.Address.Hex(), amountToSend.String(), gasCost.String(), nonce)
	s.logger.Info("Transaction details",
		"sessionKeyAddress", sessionKeyAddr.Hex(),
		"recipientAddress", firstAccount.Address.Hex(),
		"totalBalance", balance.String(),
		"gasPrice", gasPrice.String(),
		"gasLimit", gasLimit,
		"gasCost", gasCost.String(),
		"amountToSend", amountToSend.String(),
		"nonce", nonce)

	// Create the transfer transaction
	legacyTx := &types.LegacyTx{
		Nonce:    nonce,
		To:       firstAccount.Address,
		Value:    amountToSend,
		GasPrice: gasPrice,
		Gas:      gasLimit,
	}

	// Sign the transaction with the session key
	fmt.Printf("✍️  Signing transaction with session key %s\n", sessionKeyAddr.Hex())
	s.logger.Info("Signing transaction with session key", "sessionKeyAddress", sessionKeyAddr.Hex())
	signedTx, err := s.services.SKManager.SignTx(ctx, user, sessionKeyAddr, types.NewTx(legacyTx))
	if err != nil {
		fmt.Printf("❌ Failed to sign transaction: %v\n", err)
		s.logger.Error("Failed to sign transaction with session key", "error", err, "sessionKeyAddress", sessionKeyAddr.Hex())
		return fmt.Errorf("failed to sign transaction with session key: %w", err)
	}
	fmt.Printf("✅ Transaction signed successfully\n")

	// Convert to raw bytes and send
	fmt.Printf("📦 Marshaling signed transaction...\n")
	s.logger.Info("Marshaling signed transaction")
	blob, err := signedTx.MarshalBinary()
	if err != nil {
		fmt.Printf("❌ Failed to marshal transaction: %v\n", err)
		s.logger.Error("Failed to marshal signed transaction", "error", err)
		return fmt.Errorf("failed to marshal signed transaction: %w", err)
	}
	fmt.Printf("📦 Transaction marshaled, size: %d bytes\n", len(blob))

	// Send the transaction using the same approach as SendRawTx
	fmt.Printf("📤 Sending raw transaction...\n")
	s.logger.Info("Sending raw transaction", "blobSize", len(blob))
	hash, err := s.sendRawTransaction(ctx, blob, user, sessionKeyAddr)
	if err != nil {
		fmt.Printf("❌ Failed to send raw transaction: %v\n", err)
		s.logger.Error("Failed to send transaction", "error", err, "sessionKeyAddress", sessionKeyAddr.Hex())
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("🎉 SUCCESS! Fund recovery transaction sent - Hash: %s\n", hash.Hex())
	s.logger.Info("Successfully sent fund recovery transaction",
		"userID", wecommon.HashForLogging(user.ID),
		"sessionKeyAddress", sessionKeyAddr.Hex(),
		"recipientAddress", firstAccount.Address.Hex(),
		"amountToSend", amountToSend.String(),
		"gasCost", gasCost.String(),
		"txHash", hash.Hex())

	return nil
}

// getGasPrice retrieves the current gas price
func (s *SessionKeyExpirationService) getGasPrice(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	_, err := WithPlainRPCConnection(ctx, s.services.BackendRPC, func(client *rpc.Client) (*hexutil.Big, error) {
		err := client.CallContext(ctx, &result, "ten_gasPrice")
		return &result, err
	})
	if err != nil {
		s.logger.Error("Failed to get gas price", "error", err)
		return nil, fmt.Errorf("failed to get gas price via RPC: %w", err)
	}
	gasPrice := result.ToInt()
	s.logger.Info("Retrieved gas price", "gasPrice", gasPrice.String())
	return gasPrice, nil
}

// estimateGas estimates gas for a transfer transaction
func (s *SessionKeyExpirationService) estimateGas(ctx context.Context, from, to common.Address, value *hexutil.Big) (uint64, error) {
	var result hexutil.Uint64
	_, err := WithPlainRPCConnection(ctx, s.services.BackendRPC, func(client *rpc.Client) (*hexutil.Uint64, error) {
		err := client.CallContext(ctx, &result, "ten_estimateGas", map[string]interface{}{
			"from":  from.Hex(),
			"to":    to.Hex(),
			"value": value.String(),
		})
		return &result, err
	})
	if err != nil {
		s.logger.Error("Failed to estimate gas", "error", err, "from", from.Hex(), "to", to.Hex(), "value", value.String())
		return 0, fmt.Errorf("failed to estimate gas via RPC: %w", err)
	}
	gasLimit := uint64(result)
	s.logger.Info("Estimated gas", "gasLimit", gasLimit, "from", from.Hex(), "to", to.Hex(), "value", value.String())
	return gasLimit, nil
}

// getNonce retrieves the nonce for an address
func (s *SessionKeyExpirationService) getNonce(ctx context.Context, addr common.Address) (uint64, error) {
	var result hexutil.Uint64
	_, err := WithPlainRPCConnection(ctx, s.services.BackendRPC, func(client *rpc.Client) (*hexutil.Uint64, error) {
		err := client.CallContext(ctx, &result, "ten_getTransactionCount", addr, rpc.LatestBlockNumber)
		return &result, err
	})
	if err != nil {
		s.logger.Error("Failed to get nonce", "error", err, "address", addr.Hex())
		return 0, fmt.Errorf("failed to get nonce via RPC: %w", err)
	}
	nonce := uint64(result)
	s.logger.Info("Retrieved nonce", "nonce", nonce, "address", addr.Hex())
	return nonce, nil
}

// sendRawTransaction sends a raw transaction using the same approach as SendRawTx
func (s *SessionKeyExpirationService) sendRawTransaction(ctx context.Context, input hexutil.Bytes, user *wecommon.GWUser, sessionKeyAddr common.Address) (common.Hash, error) {
	// Get the session key account for authentication
	sessionKey, exists := user.SessionKeys[sessionKeyAddr]
	if !exists {
		return common.Hash{}, fmt.Errorf("session key not found for address %s", sessionKeyAddr.Hex())
	}

	var result common.Hash
	_, err := WithEncRPCConnection(ctx, s.services.BackendRPC, sessionKey.Account, func(rpcClient *tenrpc.EncRPCClient) (*common.Hash, error) {
		err := rpcClient.CallContext(ctx, &result, tencommonrpc.ERPCSendRawTransaction, input)
		return &result, err
	})
	if err != nil {
		return common.Hash{}, err
	}
	return result, nil
}
