package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SKManager - session keys are Private Keys managed by the Gateway
// Each user can have multiple Session Keys (up to 100), one per dApp
// Each session key is identified by its address, which serves as the session ID
// Session keys are implicitly active when used - no activation/deactivation needed
// The SK is also considered an "Account" of that user
// when the SK is created, it signs over the VK of the user so that it can interact with a node the standard way
// From the POV of the Ten network - a session key is a normal account key
type SKManager interface {
	CreateSessionKey(user *common.GWUser) (*common.GWSessionKey, error)
	DeleteSessionKey(user *common.GWUser, sessionKeyAddr gethcommon.Address) (bool, error)
	SignTx(ctx context.Context, user *common.GWUser, sessionKeyAddr gethcommon.Address, input *types.Transaction) (*types.Transaction, error)
	SetTxSender(txSender TxSender)
}

type skManager struct {
	storage         storage.UserStorage
	config          *common.Config
	logger          gethlog.Logger
	activityTracker SessionKeyActivityTracker
	txSender        TxSender
}

func NewSKManager(storage storage.UserStorage, config *common.Config, logger gethlog.Logger, tracker SessionKeyActivityTracker) SKManager {
	return &skManager{
		storage:         storage,
		config:          config,
		logger:          logger,
		activityTracker: tracker,
		txSender:        nil, // Will be set later via SetTxSender
	}
}

// CreateSessionKey - generates a fresh key and signs over the VK of the user with it
func (m *skManager) CreateSessionKey(user *common.GWUser) (*common.GWSessionKey, error) {
	// Check session key limit
	if len(user.SessionKeys) >= common.MaxSessionKeysPerUser {
		return nil, fmt.Errorf("maximum number of session keys (%d) reached", common.MaxSessionKeysPerUser)
	}

	sk, err := m.createSK(user)
	if err != nil {
		return nil, err
	}
	err = m.storage.AddSessionKey(user.ID, *sk)
	if err != nil {
		return nil, err
	}

	// Mark activity for the new session key
	if m.activityTracker != nil && sk != nil && sk.Account != nil && sk.Account.Address != nil {
		m.activityTracker.MarkActive(user.ID, *sk.Account.Address)
	}
	return sk, nil
}

// SetTxSender injects the TxSender dependency after SKManager creation
func (m *skManager) SetTxSender(txSender TxSender) {
	m.txSender = txSender
}

func (m *skManager) DeleteSessionKey(user *common.GWUser, sessionKeyAddr gethcommon.Address) (bool, error) {
	if len(user.SessionKeys) == 0 {
		return false, errors.New("no session keys found")
	}

	if _, exists := user.SessionKeys[sessionKeyAddr]; !exists {
		return false, fmt.Errorf("session key not found: %s", sessionKeyAddr.Hex())
	}

	// Before deleting the session key, attempt to refund any remaining funds
	// Find the first account registered with the user - we will send funds to this account
	var firstAccount *common.GWAccount
	for _, account := range user.Accounts {
		firstAccount = account
		break
	}

	if firstAccount != nil && firstAccount.Address != nil && m.txSender != nil {
		// Attempt to refund funds from the session key to the user's primary account
		_, err := m.txSender.SendAllMinusGasWithSK(context.Background(), user, sessionKeyAddr, *firstAccount.Address)
		if err != nil {
			m.logger.Error("Failed to refund funds from session key before deletion",
				"error", err,
				"userID", common.HashForLogging(user.ID),
				"sessionKeyAddress", sessionKeyAddr.Hex())
			// Continue with deletion even if refund fails - don't block the deletion
		} else {
			m.logger.Info("Successfully refunded funds from session key before deletion",
				"userID", common.HashForLogging(user.ID),
				"sessionKeyAddress", sessionKeyAddr.Hex())
		}
	}

	// Remove from activity tracker before deleting from storage
	if m.activityTracker != nil {
		_ = m.activityTracker.Delete(sessionKeyAddr)
	}

	err := m.storage.RemoveSessionKey(user.ID, &sessionKeyAddr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *skManager) createSK(user *common.GWUser) (*common.GWSessionKey, error) {
	// generate new key-pair
	sk, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key-pair: %w", err)
	}
	skEcies := ecies.ImportECDSA(sk)

	// Compute the Ethereum address from the public key
	address := crypto.PubkeyToAddress(sk.PublicKey)

	// use the viewing key to sign over the session key
	msg, err := viewingkey.GenerateMessage(user.ID, int64(m.config.TenChainID), 1, viewingkey.EIP712Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message. Cause %w", err)
	}

	msgHash, err := viewingkey.GetMessageHash(msg, viewingkey.EIP712Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message hash. Cause %w", err)
	}

	// current signature is valid - return account address
	sig, err := crypto.Sign(msgHash, sk)
	if err != nil {
		return nil, fmt.Errorf("cannot sign message with session key. Cause %w", err)
	}

	return &common.GWSessionKey{
		PrivateKey: skEcies,
		Account: &common.GWAccount{
			User:          user,
			Address:       &address,
			Signature:     sig,
			SignatureType: viewingkey.EIP712Signature,
		},
		CreatedAt: time.Now(),
	}, nil
}

func (m *skManager) GetSessionKey(user *common.GWUser, sessionKeyAddr gethcommon.Address) (*common.GWSessionKey, error) {
	if user.SessionKeys == nil {
		return nil, errors.New("no session keys found")
	}

	sessionKey, exists := user.SessionKeys[sessionKeyAddr]
	if !exists {
		return nil, fmt.Errorf("session key not found: %s", sessionKeyAddr.Hex())
	}

	return sessionKey, nil
}

func (m *skManager) SignTx(ctx context.Context, user *common.GWUser, sessionKeyAddr gethcommon.Address, tx *types.Transaction) (*types.Transaction, error) {
	sessionKey, err := m.GetSessionKey(user, sessionKeyAddr)
	if err != nil {
		return nil, err
	}

	prvKey := sessionKey.PrivateKey.ExportECDSA()
	signer := types.NewCancunSigner(big.NewInt(int64(m.config.TenChainID)))

	stx, err := types.SignTx(tx, signer, prvKey)
	if err != nil {
		return nil, err
	}

	m.logger.Debug("Signed transaction with session key", "stxHash", stx.Hash().Hex(), "sessionKey", sessionKeyAddr.Hex())
	if m.activityTracker != nil {
		m.activityTracker.MarkActive(user.ID, sessionKeyAddr)
	}

	return stx, nil
}
