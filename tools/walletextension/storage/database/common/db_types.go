package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	common2 "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type GWUserDB struct {
	UserId     []byte        `json:"userId"`
	PrivateKey []byte        `json:"privateKey"`
	Accounts   []GWAccountDB `json:"accounts"`
}

type GWAccountDB struct {
	AccountAddress []byte `json:"accountAddress"`
	Signature      []byte `json:"signature"`
	SignatureType  int    `json:"signatureType"`
}

func (userDB *GWUserDB) ToGWUser() *common2.GWUser {
	result := &common2.GWUser{
		UserID:   userDB.UserId,
		Accounts: make(map[common.Address]*common2.GWAccount),
		UserKey:  userDB.PrivateKey,
	}

	for _, accountDB := range userDB.Accounts {
		address := common.BytesToAddress(accountDB.AccountAddress)
		gwAccount := &common2.GWAccount{
			User:          result,
			Address:       &address,
			Signature:     accountDB.Signature,
			SignatureType: viewingkey.SignatureType(accountDB.SignatureType),
		}
		result.Accounts[address] = gwAccount
	}

	return result
}
