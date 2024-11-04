package services

import (
	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ethereum/go-ethereum/common"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

var userCacheKeyPrefix = []byte{0x0, 0x1, 0x2, 0x3}

type GWAccount struct {
	user          *GWUser
	Address       *common.Address
	signature     []byte
	signatureType viewingkey.SignatureType
}

type GWUser struct {
	userID   []byte
	services *Services
	Accounts map[common.Address]*GWAccount
	userKey  []byte
}

func (u GWUser) GetAllAddresses() []*common.Address {
	accts := make([]*common.Address, 0)
	for _, acc := range u.Accounts {
		accts = append(accts, acc.Address)
	}
	return accts
}

func gwUserFromDB(userDB wecommon.GWUserDB, s *Services) (*GWUser, error) {
	result := &GWUser{
		userID:   userDB.UserId,
		services: s,
		Accounts: make(map[common.Address]*GWAccount),
		userKey:  userDB.PrivateKey,
	}

	for _, accountDB := range userDB.Accounts {
		address := common.BytesToAddress(accountDB.AccountAddress)
		gwAccount := &GWAccount{
			user:          result,
			Address:       &address,
			signature:     accountDB.Signature,
			signatureType: viewingkey.SignatureType(accountDB.SignatureType),
		}
		result.Accounts[address] = gwAccount
	}

	return result, nil
}

func userCacheKey(userID []byte) []byte {
	var key []byte
	key = append(key, userCacheKeyPrefix...)
	key = append(key, userID...)
	return key
}
