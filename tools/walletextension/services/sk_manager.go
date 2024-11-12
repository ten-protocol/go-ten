package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SKManager - session keys are Private Keys managed by the Gateway
// At the moment, each user can have a single Session Key. Which is either active or inactive
// when the SK is active, then all transactions submitted by that user will be signed with the session key
// The SK is also considered an "Account" of that user
// when the SK is created, it signs over the VK of the user so that it can interact with a node the standard way
// From the POV of the Ten network - a session key is a normal account key
type SKManager interface {
	CreateSessionKey(user *common.GWUser) (*common.GWSessionKey, error)
	SignTx(ctx context.Context, user *common.GWUser, input *types.Transaction) (*types.Transaction, error)
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

func (m *skManager) SignTx(ctx context.Context, user *common.GWUser, tx *types.Transaction) (*types.Transaction, error) {
	prvKey := user.SessionKey.PrivateKey.ExportECDSA()
	signer := types.NewCancunSigner(big.NewInt(int64(m.config.TenChainID)))

	stx, err := types.SignTx(tx, signer, prvKey)
	if err != nil {
		return nil, err
	}

	m.logger.Debug("Signed transaction with session key", "stxHash", stx.Hash().Hex())

	return stx, nil
}
