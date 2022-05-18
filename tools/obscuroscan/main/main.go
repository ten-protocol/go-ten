package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/tools/obscuroscan"
	"strconv"
)

const localhost = "localhost:"

func main() {
	config := parseCLIArgs()
	nodeID := common.BytesToAddress([]byte(config.nodeID))
	obscuroscanAddr := localhost + strconv.Itoa(config.startPort)

	server := obscuroscan.NewObscuroscan(nodeID, config.clientServerAddr)
	go server.Serve(obscuroscanAddr)
	fmt.Printf("Obscuroscan started.\n💡 Visit %s to monitor the Obscuro network.\n", obscuroscanAddr)

	defer server.Shutdown()
	select {}
}
