package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
)

func main() {
	config := parseCLIArgs()

	nodeAddress := common.BytesToAddress([]byte(*config.nodeID))
	if err := enclave.StartServer(*config.port, nodeAddress, nil); err != nil {
		panic(err)
	}
	fmt.Printf("Enclave server listening on port %d.\n", *config.port)

	select {}
}
