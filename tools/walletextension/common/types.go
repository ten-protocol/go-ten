package common

import (
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

type SessionKeyActivity struct {
	Addr       common.Address
	UserID     []byte
	LastActive time.Time
}
