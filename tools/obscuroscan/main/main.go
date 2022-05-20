package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/tools/obscuroscan"
)

func main() {
	config := parseCLIArgs()
	nodeID := common.BytesToAddress([]byte(config.nodeID))

	server := obscuroscan.NewObscuroscan(nodeID, config.clientServerAddr)
	go server.Serve(config.address)
	fmt.Printf("Obscuroscan started.\n💡 Visit %s to monitor the Obscuro network.\n", config.address)

	defer server.Shutdown()
	select {}
}
