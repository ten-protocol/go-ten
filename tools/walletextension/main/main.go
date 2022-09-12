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
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Welcome to the Obscuro wallet extension. \n\n")
	fmt.Printf("Starting with following config: \n%s\n", string(jsonConfig))

	// We wait thirty seconds for a connection to the node. If we cannot establish one, we exit the program.
	fmt.Printf("Waiting up to thirty seconds for connection to host at %s...\n", config.NodeRPCWebsocketAddress)
	counter := 30
	for {
		conn, err := net.Dial(tcp, config.NodeRPCWebsocketAddress)
		if conn != nil {
			conn.Close()
		}
		if err == nil {
			break
		}

		counter--
		if counter <= 0 {
			fmt.Printf("Exiting. Could not establish connection to host at %s. Cause: %s\n", config.NodeRPCWebsocketAddress, err)
			return
		}
		time.Sleep(time.Second)
	}

	walletExtension := walletextension.NewWalletExtension(config)
	defer walletExtension.Shutdown()

	go walletExtension.Serve(localhost, config.WalletExtensionPort, config.WalletExtensionPortWS)

	walletExtensionAddr := fmt.Sprintf("%s:%d", localhost, config.WalletExtensionPort)
	fmt.Printf("💡 Wallet extension started - visit http://%s/viewingkeys/ to generate an ephemeral viewing key.\n", walletExtensionAddr)

	select {}
}
