package rpcapi

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/rpc"
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

func userCacheKey(userID []byte) []byte {
	var key []byte
	key = append(key, userCacheKeyPrefix...)
	key = append(key, userID...)
	return key
}

func getUser(userID []byte, s *Services) (*GWUser, error) {
	return withCache(s.Cache, &CacheCfg{TTL: longCacheTTL}, userCacheKey(userID), func() (*GWUser, error) {
		result := GWUser{userID: userID, accounts: map[common.Address]*GWAccount{}}
		userPrivateKey, err := s.Storage.GetUserPrivateKey(userID)
		if err != nil {
			return nil, fmt.Errorf("user %s not found. %w", hexutils.BytesToHex(userID), err)
		}
		result.userKey = userPrivateKey
		allAccounts, err := s.Storage.GetAccounts(userID)
		if err != nil {
			return nil, err
		}

		for _, account := range allAccounts {
			address := common.BytesToAddress(account.AccountAddress)
			result.accounts[address] = &GWAccount{user: &result, address: &address, signature: account.Signature, signatureType: viewingkey.SignatureType(uint8(account.SignatureType))}
		}
		return &result, nil
	})
}

func (account *GWAccount) connect(url string, logger gethlog.Logger) (*rpc.EncRPCClient, error) {
	// create a new client
	// todo - close and cache
	encClient, err := wecommon.CreateEncClient(url, account.address.Bytes(), account.user.userKey, account.signature, account.signatureType, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating new client, %w", err)
	}
	return encClient, nil
}
