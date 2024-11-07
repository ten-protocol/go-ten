package services

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

type SKManager interface {
	CreateSessionKey(user *common.GWUser) (*common.GWSessionKey, error)
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
	err = m.storage.AddSessionKey(user.UserID, *sk)
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
	msg, err := viewingkey.GenerateMessage(user.UserID, int64(m.config.TenChainID), 1, viewingkey.EIP712Signature)
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
