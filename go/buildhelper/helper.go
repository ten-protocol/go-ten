package buildhelper

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/buildhelper/buildconstants"
)

func Addr1PK() *ecdsa.PrivateKey {
	privateKey, _ := crypto.HexToECDSA(buildconstants.ADDR1PK)
	return privateKey
}

func Addr2PK() *ecdsa.PrivateKey {
	privateKey, _ := crypto.HexToECDSA(buildconstants.ADDR2PK)
	return privateKey
}

func Addr3PK() *ecdsa.PrivateKey {
	privateKey, _ := crypto.HexToECDSA(buildconstants.ADDR3PK)
	return privateKey
}

func Addr1() common.Address {
	privateKey := Addr1PK()
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func Addr2() common.Address {
	privateKey := Addr2PK()
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func Addr3() common.Address {
	privateKey := Addr3PK()
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
