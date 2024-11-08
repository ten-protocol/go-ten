package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
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

func (userDB *GWUserDB) ToGWUser() *wecommon.GWUser {
	result := &wecommon.GWUser{
		UserID:   userDB.UserId,
		Accounts: make(map[common.Address]*wecommon.GWAccount),
		UserKey:  userDB.PrivateKey,
	}

	for _, accountDB := range userDB.Accounts {
		address := common.BytesToAddress(accountDB.AccountAddress)
		gwAccount := wecommon.GWAccount{
			User:          result,
			Address:       &address,
			Signature:     accountDB.Signature,
			SignatureType: viewingkey.SignatureType(accountDB.SignatureType),
		}
		result.Accounts[address] = &gwAccount
	}

	return result
}
