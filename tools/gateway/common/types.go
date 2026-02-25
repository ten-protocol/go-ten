package common

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"golang.org/x/exp/maps"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// MaxSessionKeysPerUser defines the maximum number of session keys a user can have
	MaxSessionKeysPerUser = 100
)

// GWSessionKey - an account key-pair registered for a user
type GWSessionKey struct {
	Account    *GWAccount
	PrivateKey *ecies.PrivateKey // the private key corresponding to the account
	CreatedAt  time.Time         // timestamp when the session key was created
}

type GWAccount struct {
	User          *GWUser
	Address       *common.Address
	Signature     []byte // the signature by the account over the userId - which is derived from the VK
	SignatureType viewingkey.SignatureType
}

type GWUser struct {
	ID          []byte
	Accounts    map[common.Address]*GWAccount
	UserKey     []byte
	SessionKeys map[common.Address]*GWSessionKey // map of session key address to session key
}

func (u GWUser) AllAccounts() map[common.Address]*GWAccount {
	res := maps.Clone(u.Accounts)
	for addr, sessionKey := range u.SessionKeys {
		res[addr] = sessionKey.Account
	}
	return res
}

func (u GWUser) GetAllAddresses() []common.Address {
	return maps.Keys(u.AllAccounts())
}

// GetFirstAccount returns the first account from the user's Accounts map.
// Returns an error if the user has no accounts or if the account has no address.
func (u *GWUser) GetFirstAccount() (*GWAccount, error) {
	if len(u.Accounts) == 0 {
		return nil, fmt.Errorf("user has no accounts")
	}

	for _, account := range u.Accounts {
		if account == nil {
			continue
		}
		if account.Address == nil {
			return nil, fmt.Errorf("account has no address")
		}
		return account, nil
	}

	return nil, fmt.Errorf("no valid account found")
}

type SessionKeyActivity struct {
	Addr       common.Address
	UserID     []byte
	LastActive time.Time
}
