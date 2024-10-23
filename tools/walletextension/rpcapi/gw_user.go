package rpcapi

import (
	"fmt"

	"github.com/status-im/keycard-go/hexutils"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ethereum/go-ethereum/common"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

var userCacheKeyPrefix = []byte{0x0, 0x1, 0x2, 0x3}

type GWAccount struct {
	user          *GWUser
	address       *common.Address
	signature     []byte
	signatureType viewingkey.SignatureType
}

type GWUser struct {
	userID   []byte
	services *Services
	accounts map[common.Address]*GWAccount
	userKey  []byte
}

func (u GWUser) GetAllAddresses() []*common.Address {
	accts := make([]*common.Address, 0)
	for _, acc := range u.accounts {
		accts = append(accts, acc.address)
	}
	return accts
}

func gwUserFromDB(userDB wecommon.GWUserDB, s *Services) (*GWUser, error) {
	result := &GWUser{
		userID:   userDB.UserId,
		services: s,
		accounts: make(map[common.Address]*GWAccount),
		userKey:  userDB.PrivateKey,
	}

	for _, accountDB := range userDB.Accounts {
		address := common.BytesToAddress(accountDB.AccountAddress)
		gwAccount := &GWAccount{
			user:          result,
			address:       &address,
			signature:     accountDB.Signature,
			signatureType: viewingkey.SignatureType(accountDB.SignatureType),
		}
		result.accounts[address] = gwAccount
	}

	return result, nil
}

func userCacheKey(userID []byte) []byte {
	var key []byte
	key = append(key, userCacheKeyPrefix...)
	key = append(key, userID...)
	return key
}

func getUser(userID []byte, s *Services) (*GWUser, error) {
	return withCache(s.Cache, &CacheCfg{CacheType: LongLiving}, userCacheKey(userID), func() (*GWUser, error) {
		user, err := s.Storage.GetUser(userID)
		if err != nil {
			return nil, fmt.Errorf("user %s not found. %w", hexutils.BytesToHex(userID), err)
		}
		result, err := gwUserFromDB(user, s)
		return result, err
	})
}
