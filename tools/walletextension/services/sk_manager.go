package services

import (
	"context"
	"fmt"
	"math/big"

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
	ListSessionKeys(user *common.GWUser) ([]gethcommon.Address, error)
	GetSessionKey(user *common.GWUser, sessionKeyAddr gethcommon.Address) (*common.GWSessionKey, error)
	SignTx(ctx context.Context, user *common.GWUser, sessionKeyAddr gethcommon.Address, input *types.Transaction) (*types.Transaction, error)
}

type skManager struct {
	storage storage.UserStorage
	config  *common.Config
	logger  gethlog.Logger
}

func NewSKManager(storage storage.UserStorage, config *common.Config, logger gethlog.Logger) SKManager {
	return &skManager{
		storage: storage,
		config:  config,
		logger:  logger,
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
	return sk, nil
}


func (m *skManager) DeleteSessionKey(user *common.GWUser, sessionKeyAddr gethcommon.Address) (bool, error) {
	if user.SessionKeys == nil || len(user.SessionKeys) == 0 {
		return false, fmt.Errorf("no session keys found")
	}

	if _, exists := user.SessionKeys[sessionKeyAddr]; !exists {
		return false, fmt.Errorf("session key not found: %s", sessionKeyAddr.Hex())
	}

	err := m.storage.RemoveSessionKey(user.ID, sessionKeyAddr.Bytes())
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
	}, nil
}

func (m *skManager) ListSessionKeys(user *common.GWUser) ([]gethcommon.Address, error) {
	if user.SessionKeys == nil || len(user.SessionKeys) == 0 {
		return []gethcommon.Address{}, nil
	}

	addresses := make([]gethcommon.Address, 0, len(user.SessionKeys))
	for addr := range user.SessionKeys {
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

func (m *skManager) GetSessionKey(user *common.GWUser, sessionKeyAddr gethcommon.Address) (*common.GWSessionKey, error) {
	if user.SessionKeys == nil {
		return nil, fmt.Errorf("no session keys found")
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

	return stx, nil
}
