package common

import (
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"golang.org/x/exp/maps"

	"github.com/ethereum/go-ethereum/common"
)

type GWAccount struct {
	User          *GWUser
	Address       *common.Address
	Signature     []byte
	SignatureType viewingkey.SignatureType
}

type GWUser struct {
	UserID   []byte
	Accounts map[common.Address]*GWAccount
	UserKey  []byte
}

func (u GWUser) GetAllAddresses() []common.Address {
	return maps.Keys(u.Accounts)
}
