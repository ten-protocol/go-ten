package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKey, _ := crypto.GenerateKey()
	hex := common.Bytes2Hex(crypto.FromECDSA(privateKey))
	println(hex)

	hex2 := common.Bytes2Hex(crypto.CompressPubkey(&privateKey.PublicKey))
	println(hex2)
}
