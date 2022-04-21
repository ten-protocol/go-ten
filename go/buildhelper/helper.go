package buildhelper

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/buildconstants"
)

func Addr1PK() *ecdsa.PrivateKey {
	privateKey, _ := crypto.HexToECDSA(buildconstants.Addr1PK)
	return privateKey
}

func Addr2PK() *ecdsa.PrivateKey {
	privateKey, _ := crypto.HexToECDSA(buildconstants.Addr2PK)
	return privateKey
}

func Addr3PK() *ecdsa.PrivateKey {
	privateKey, _ := crypto.HexToECDSA(buildconstants.Addr3PK)
	return privateKey
}

func Addr1() common.Address {
	privateKey := Addr1PK()
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func Addr2() common.Address {
	privateKey := Addr2PK()
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func Addr3() common.Address {
	privateKey := Addr3PK()
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
