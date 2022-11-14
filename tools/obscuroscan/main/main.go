package main

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/tools/obscuroscan"
)

func main() {
	config := parseCLIArgs()

	server := obscuroscan.NewObscuroscan(
		config.rpcServerAddr,
		log.New(log.ObscuroscanCmp, int(gethlog.LvlInfo), config.logPath),
	)
	go server.Serve(config.address)
	fmt.Printf("Obscuroscan started.\n💡 Visit %s to monitor the Obscuro network.\n", config.address)

	defer server.Shutdown()
	select {}
}
