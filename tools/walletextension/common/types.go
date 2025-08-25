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

	// Session key expiry constants
	SessionKeyExpiryDuration = 24 * time.Hour  // Session keys expire after 24 hours
	MinETHReturnThreshold    = 0.001           // Minimum ETH balance to trigger return (in ETH)
	ExpiryCheckInterval      = 1 * time.Hour   // How often to check for expired session keys
)

// GWSessionKey - an account key-pair registered for a user
type GWSessionKey struct {
	Account    *GWAccount
	PrivateKey *ecies.PrivateKey // the private key corresponding to the account
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
