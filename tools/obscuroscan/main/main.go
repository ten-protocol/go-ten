package main

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/tools/obscuroscan"
)

func main() {
	config := parseCLIArgs()

	server := obscuroscan.NewObscuroscan(config.clientServerAddr)
	go server.Serve(config.address)
	fmt.Printf("Obscuroscan started.\n💡 Visit %s to monitor the Obscuro network.\n", config.address)

	defer server.Shutdown()
	select {}
}
