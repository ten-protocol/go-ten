package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"math/big"
)

func main() {
	client := obscuroclient.NewClient(common.BigToAddress(big.NewInt(1)), "20.68.160.65:13000")
	err := client.Call(nil, obscuroclient.RPCGetCurrentBlockHeadHeight)
	if err != nil {
		panic(err)
	}
}
