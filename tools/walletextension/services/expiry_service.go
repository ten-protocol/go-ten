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
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	tenrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SessionKeyMaintenanceService manages recurring session key fund returns
type SessionKeyMaintenanceService interface {
	// ProcessRecurringFundReturns finds and processes session keys due for fund return
	ProcessRecurringFundReturns(ctx context.Context) error
	// Start begins the background fund return checking process
	Start(ctx context.Context)
	// Stop gracefully stops the maintenance service
	Stop()
}

type sessionKeyMaintenanceService struct {
	storage   storage.UserStorage
	rpc       BackendRPC
	skManager SKManager
	logger    gethlog.Logger
	stopChan  chan struct{}
}

// NewSessionKeyMaintenanceService creates a new session key maintenance service
func NewSessionKeyMaintenanceService(storage storage.UserStorage, rpc *BackendRPC, skManager SKManager, logger gethlog.Logger) SessionKeyMaintenanceService {
	return &sessionKeyMaintenanceService{
		storage:   storage,
		rpc:       *rpc,
		skManager: skManager,
		logger:    logger,
		stopChan:  make(chan struct{}),
	}
}

// Start begins the background fund return checking process
func (s *sessionKeyMaintenanceService) Start(ctx context.Context) {
	ticker := time.NewTicker(wecommon.ExpiryCheckInterval)
	defer ticker.Stop()

	s.logger.Info("Starting session key maintenance service for recurring fund returns", "interval", wecommon.ExpiryCheckInterval)

	for {
		select {
		case <-ticker.C:
			if err := s.ProcessRecurringFundReturns(ctx); err != nil {
				s.logger.Error("Failed to process recurring fund returns for session keys", "error", err)
			}
		case <-s.stopChan:
			s.logger.Info("Session key maintenance service stopped")
			return
		case <-ctx.Done():
			s.logger.Info("Session key maintenance service context canceled")
			return
		}
	}
}

// Stop gracefully stops the maintenance service
func (s *sessionKeyMaintenanceService) Stop() {
	close(s.stopChan)
}

// ProcessRecurringFundReturns finds and processes session keys due for fund return
func (s *sessionKeyMaintenanceService) ProcessRecurringFundReturns(ctx context.Context) error {
	// This is a simplified implementation - in a production system you'd want to
	// iterate through users more efficiently, possibly with pagination
	// For now, we'll need a way to get all users from storage
	// This would require extending the storage interface
	
	s.logger.Debug("Starting recurring fund return processing for all session keys")
	
	// TODO: Implement user iteration to check all session keys for recurring fund returns
	// This requires extending the storage interface to iterate through users
	
	return nil
}

// processSessionKeyFundReturnsForUser handles recurring fund returns for a specific user
func (s *sessionKeyMaintenanceService) processSessionKeyFundReturnsForUser(ctx context.Context, user *wecommon.GWUser) error {
	now := time.Now().Unix()
	
	// Get user's first account (primary account to receive returned funds)
	var primaryAccount *wecommon.GWAccount
	for _, account := range user.Accounts {
		primaryAccount = account
		break // Take the first account
	}
	
	if primaryAccount == nil {
		return fmt.Errorf("user has no primary account to return funds to")
	}
	
	// Check each session key for recurring fund return (every 24 hours)
	for sessionKeyAddr, sessionKey := range user.SessionKeys {
		// Check if 24 hours have passed since last fund return (recurring cycle)
		if s.shouldReturnFunds(sessionKeyAddr, now) {
			if err := s.returnSessionKeyFunds(ctx, sessionKey, primaryAccount, sessionKeyAddr); err != nil {
				s.logger.Error("Failed to return funds for session key", 
					"sessionKey", sessionKeyAddr.Hex(), "error", err)
				continue
			}
		}
	}
	
	return nil
}

// shouldReturnFunds checks if 24 hours have passed since last fund return
func (s *sessionKeyMaintenanceService) shouldReturnFunds(sessionKeyAddr common.Address, now int64) bool {
	// TODO: Query database for LastFundReturn timestamp and check if 24 hours have passed
	// This requires extending the storage interface to access session key metadata
	// For now, return false to prevent fund returns until properly implemented
	return false
}

// returnSessionKeyFunds transfers remaining ETH from session key back to primary account (recurring 24h cycle) and updates LastFundReturn
func (s *sessionKeyMaintenanceService) returnSessionKeyFunds(ctx context.Context, sessionKey *wecommon.GWSessionKey, primaryAccount *wecommon.GWAccount, sessionKeyAddr common.Address) error {
	// Check ETH balance of session key
	balance, err := s.getETHBalance(ctx, *sessionKey.Account.Address)
	if err != nil {
		return fmt.Errorf("failed to get session key balance: %w", err)
	}
	
	// Convert threshold to wei (0.001 ETH = 1e15 wei)
	thresholdWei := new(big.Int).Mul(big.NewInt(int64(wecommon.MinETHReturnThreshold*1000)), big.NewInt(1e15))
	
	// Skip if balance is below threshold
	if balance.Cmp(thresholdWei) < 0 {
		s.logger.Debug("Session key balance below threshold, skipping", 
			"balance", balance.String(), "threshold", thresholdWei.String())
		return nil
	}
	
	// Estimate gas for transfer
	gasLimit := uint64(21000) // Standard ETH transfer gas limit
	gasPrice, err := s.getGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}
	
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	
	// Calculate amount to transfer (balance - gas cost)
	transferAmount := new(big.Int).Sub(balance, gasCost)
	if transferAmount.Sign() <= 0 {
		s.logger.Debug("Insufficient balance to cover gas fees", 
			"balance", balance.String(), "gasCost", gasCost.String())
		return nil
	}
	
	// Create and sign transfer transaction
	nonce, err := s.getNonce(ctx, *sessionKey.Account.Address)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	
	tx := types.NewTransaction(nonce, *primaryAccount.Address, transferAmount, gasLimit, gasPrice, nil)
	
	// Sign transaction with session key
	signedTx, err := s.skManager.SignTx(ctx, sessionKey.Account.User, *sessionKey.Account.Address, tx)
	if err != nil {
		return fmt.Errorf("failed to sign transfer transaction: %w", err)
	}
	
	// Submit transaction
	if err := s.submitTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to submit transfer transaction: %w", err)
	}
	
	s.logger.Info("Successfully submitted fund return transaction", 
		"from", sessionKey.Account.Address.Hex(), 
		"to", primaryAccount.Address.Hex(), 
		"amount", transferAmount.String(),
		"txHash", signedTx.Hash().Hex())
	
	// Update LastFundReturn timestamp to start next 24-hour recurring cycle
	if err := s.updateLastFundReturn(sessionKeyAddr, time.Now().Unix()); err != nil {
		s.logger.Error("Failed to update LastFundReturn timestamp", 
			"sessionKey", sessionKeyAddr.Hex(), "error", err)
		// Don't return error here - the fund transfer succeeded, next cycle will continue
	}
	
	return nil
}

// getETHBalance retrieves the ETH balance for an address
func (s *sessionKeyMaintenanceService) getETHBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return WithPlainRPCConnection(ctx, &s.rpc, func(client *gethrpc.Client) (*big.Int, error) {
		var result hexutil.Big
		err := client.CallContext(ctx, &result, tenrpc.ERPCGetBalance, address, gethrpc.LatestBlockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get balance for %s: %w", address.Hex(), err)
		}
		return (*big.Int)(&result), nil
	})
}

// getGasPrice retrieves the current gas price
func (s *sessionKeyMaintenanceService) getGasPrice(ctx context.Context) (*big.Int, error) {
	return WithPlainRPCConnection(ctx, &s.rpc, func(client *gethrpc.Client) (*big.Int, error) {
		var result hexutil.Big
		err := client.CallContext(ctx, &result, "eth_gasPrice")
		if err != nil {
			s.logger.Warn("Failed to get gas price, using fallback", "error", err)
			return big.NewInt(20000000000), nil // 20 gwei fallback
		}
		return (*big.Int)(&result), nil
	})
}

// getNonce retrieves the current nonce for an address
func (s *sessionKeyMaintenanceService) getNonce(ctx context.Context, address common.Address) (uint64, error) {
	result, err := WithPlainRPCConnection(ctx, &s.rpc, func(client *gethrpc.Client) (*hexutil.Uint64, error) {
		var result hexutil.Uint64
		err := client.CallContext(ctx, &result, tenrpc.ERPCGetTransactionCount, address, gethrpc.LatestBlockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce for %s: %w", address.Hex(), err)
		}
		return &result, nil
	})
	if err != nil {
		return 0, err
	}
	return uint64(*result), nil
}

// submitTransaction submits a signed transaction to the network
func (s *sessionKeyMaintenanceService) submitTransaction(ctx context.Context, tx *types.Transaction) error {
	blob, err := tx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	result, err := WithPlainRPCConnection(ctx, &s.rpc, func(client *gethrpc.Client) (*common.Hash, error) {
		var result common.Hash
		err := client.CallContext(ctx, &result, tenrpc.ERPCSendRawTransaction, hexutil.Encode(blob))
		if err != nil {
			return nil, fmt.Errorf("failed to send transaction: %w", err)
		}
		return &result, nil
	})
	if err != nil {
		return err
	}
	s.logger.Info("Transaction submitted successfully", "hash", result.Hex())
	return nil
}

// updateLastFundReturn updates the LastFundReturn timestamp to start the next recurring 24-hour cycle
func (s *sessionKeyMaintenanceService) updateLastFundReturn(sessionKeyAddr common.Address, timestamp int64) error {
	// TODO: Extend storage interface to update session key metadata
	// This update starts the next 24-hour recurring fund return cycle
	s.logger.Debug("Would update LastFundReturn for next recurring cycle", 
		"sessionKey", sessionKeyAddr.Hex(), "nextCycleStartsAt", timestamp)
	return nil
}