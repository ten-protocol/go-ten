package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	tcp       = "tcp"
	localhost = "127.0.0.1"
)

func main() {
	config := parseCLIArgs()

	// We wait thirty seconds for a connection to the node. If we cannot establish one, we exit the program.
	fmt.Printf("Waiting up to thirty seconds for connection to host at %s...\n", config.NodeRPCHTTPAddress)
	counter := 30
	for {
		conn, err := net.Dial(tcp, config.NodeRPCHTTPAddress)
		if conn != nil {
			conn.Close()
		}
		if err == nil {
			break
		}

		counter--
		if counter <= 0 {
			fmt.Printf("Could not establish connection to host at %s. Exiting.", config.NodeRPCHTTPAddress)
			return
		}
		time.Sleep(time.Second)
	}

	walletExtension := walletextension.NewWalletExtension(config)
	defer walletExtension.Shutdown()
	walletExtensionAddr := fmt.Sprintf("%s:%d", localhost, config.WalletExtensionPort)
	go walletExtension.Serve(walletExtensionAddr)
	s, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Wallet extension started with following config: \n%s\n\n", string(s))
	fmt.Printf("💡 Visit http://%s/viewingkeys/ to generate an ephemeral viewing key.\n", walletExtensionAddr)

	select {}
}
