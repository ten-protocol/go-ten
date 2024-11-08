package common

import (
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"golang.org/x/exp/maps"

	"github.com/ethereum/go-ethereum/common"
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
	UserID     []byte
	Accounts   map[common.Address]*GWAccount
	UserKey    []byte
	SessionKey *GWSessionKey
	ActiveSK   bool // the session key is active, and it must be used to sign all incoming transactions, and used as the preferred account
}

func (u GWUser) AllAccounts() map[common.Address]*GWAccount {
	res := maps.Clone(u.Accounts)
	if u.SessionKey != nil {
		res[*u.SessionKey.Account.Address] = u.SessionKey.Account
	}
	return res
}

func (u GWUser) GetAllAddresses() []common.Address {
	return maps.Keys(u.AllAccounts())
}
