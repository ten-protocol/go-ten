package common

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

var ErrUserNotFound = errors.New("user not found")

type GWUserDB struct {
	UserId      []byte                             `json:"userId"`
	PrivateKey  []byte                             `json:"privateKey"`
	Accounts    []GWAccountDB                      `json:"accounts"`
	SessionKeys map[common.Address]*GWSessionKeyDB `json:"sessionKeys"` // map of session key address to session key
}

type GWAccountDB struct {
	AccountAddress []byte `json:"accountAddress"`
	Signature      []byte `json:"signature"`
	SignatureType  int    `json:"signatureType"`
}

// GWSessionKeyDB - an account key-pair registered for a user
type GWSessionKeyDB struct {
	PrivateKey     []byte      `json:"privateKey"`
	Account        GWAccountDB `json:"account"`
	CreatedAt      int64       `json:"createdAt"`      // Unix timestamp when session key was created
	LastFundReturn int64       `json:"lastFundReturn"` // Unix timestamp when funds were last returned (for recurring 24h returns)
}

func (userDB *GWUserDB) ToGWUser() (*wecommon.GWUser, error) {
	user := &wecommon.GWUser{
		ID:          userDB.UserId,
		Accounts:    make(map[common.Address]*wecommon.GWAccount),
		UserKey:     userDB.PrivateKey,
		SessionKeys: make(map[common.Address]*wecommon.GWSessionKey),
	}

	for _, accountDB := range userDB.Accounts {
		address := common.BytesToAddress(accountDB.AccountAddress)
		gwAccount := wecommon.GWAccount{
			User:          user,
			Address:       &address,
			Signature:     accountDB.Signature,
			SignatureType: viewingkey.SignatureType(accountDB.SignatureType),
		}
		user.Accounts[address] = &gwAccount
	}

	// Handle session keys
	for address, sessionKeyDB := range userDB.SessionKeys {
		ecdsaPrivateKey, err := crypto.ToECDSA(sessionKeyDB.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ECDSA private key: %w", err)
		}

		// Convert ECDSA private key to ECIES private key
		eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)
		acc := sessionKeyDB.Account
		user.SessionKeys[address] = &wecommon.GWSessionKey{
			Account: &wecommon.GWAccount{
				User:          user,
				Address:       &address,
				Signature:     acc.Signature,
				SignatureType: viewingkey.SignatureType(acc.SignatureType),
			},
			PrivateKey: eciesPrivateKey,
		}
	}

	return user, nil
}
