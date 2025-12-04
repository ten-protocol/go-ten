package services

import (
	"encoding/hex"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// ReturnDefaultUserAndAccount creates a new in-memory user and a corresponding account.
// Nothing is persisted to storage. Useful for anonymous/public flows.
func ReturnDefaultUserAndAccount(config *common.Config) (*common.GWUser, error) {
	// generate a fresh viewing key
	defaultUserVK := "5b7db1a436d96273b4ebb8a5bb28d59f28d1d54810b723dd6e03731ec335d10c" // hardcoded viewing key for the default user - remove after proper public access is implemented
	defaultUserVKBytes, err := hex.DecodeString(defaultUserVK)
	if err != nil {
		return nil, fmt.Errorf("failed to decode default VK hex: %w", err)
	}
	vk, err := crypto.ToECDSA(defaultUserVKBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert default VK bytes to private key: %w", err)
	}
	vkEcies := ecies.ImportECDSA(vk)

	// derive user ID from the viewing key
	userID := viewingkey.CalculateUserID(common.PrivateKeyToCompressedPubKey(vkEcies))

	// build an in-memory GWUser (no persistence)
	user := &common.GWUser{
		ID:          userID,
		Accounts:    make(map[gethcommon.Address]*common.GWAccount),
		UserKey:     crypto.FromECDSA(vkEcies.ExportECDSA()),
		SessionKeys: make(map[gethcommon.Address]*common.GWSessionKey),
	}

	userAddress := crypto.PubkeyToAddress(vk.PublicKey)
	msg, err := viewingkey.GenerateMessage(user.ID, int64(config.TenChainID), 1, viewingkey.EIP712Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message. Cause %w", err)
	}

	msgHash, err := viewingkey.GetMessageHash(msg, viewingkey.EIP712Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message hash. Cause %w", err)
	}

	// current signature is valid - return account address
	sig, err := crypto.Sign(msgHash, vk)
	if err != nil {
		return nil, fmt.Errorf("cannot sign message with session key. Cause %w", err)
	}

	// create an account that signs over the userID
	account := &common.GWAccount{
		User:          user,
		Address:       &userAddress,
		Signature:     sig,
		SignatureType: viewingkey.EIP712Signature,
	}

	user.Accounts[userAddress] = account

	return user, nil
}
