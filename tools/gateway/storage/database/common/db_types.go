package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	gwcommon "github.com/ten-protocol/go-ten/tools/gateway/common"
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
	PrivateKey []byte      `json:"privateKey"`
	Account    GWAccountDB `json:"account"`
	CreatedAt  time.Time   `json:"createdAt"`
}

func (userDB *GWUserDB) ToGWUser() (*gwcommon.GWUser, error) {
	user := &gwcommon.GWUser{
		ID:          userDB.UserId,
		Accounts:    make(map[common.Address]*gwcommon.GWAccount),
		UserKey:     userDB.PrivateKey,
		SessionKeys: make(map[common.Address]*gwcommon.GWSessionKey),
	}

	for _, accountDB := range userDB.Accounts {
		address := common.BytesToAddress(accountDB.AccountAddress)
		gwAccount := gwcommon.GWAccount{
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
		user.SessionKeys[address] = &gwcommon.GWSessionKey{
			Account: &gwcommon.GWAccount{
				User:          user,
				Address:       &address,
				Signature:     acc.Signature,
				SignatureType: viewingkey.SignatureType(acc.SignatureType),
			},
			PrivateKey: eciesPrivateKey,
			CreatedAt:  sessionKeyDB.CreatedAt,
		}
	}

	return user, nil
}
