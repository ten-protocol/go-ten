package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

func main() {
	config := parseCLIArgs()

	client := obscuroclient.NewClient(common.BytesToAddress([]byte(config.nodeID)), config.clientServerAddr)

	var currentBlockHeadHeight uint64
	err := client.Call(&currentBlockHeadHeight, obscuroclient.RPCGetCurrentBlockHeadHeight)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current block head height is: %d\n", currentBlockHeadHeight)
}
