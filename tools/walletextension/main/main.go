package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/container"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	tcp = "tcp"
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

	// Sets up the log file.
	if config.LogPath != log.SysOut {
		_, err := os.Create(config.LogPath)
		if err != nil {
			panic(fmt.Sprintf("could not create log file. Cause: %s", err))
		}
	}

	logLvl := gethlog.LvlError
	if config.VerboseFlag {
		logLvl = gethlog.LvlDebug
	}
	logger := log.New(log.WalletExtCmp, int(logLvl), config.LogPath)

	walletExtContainer := container.NewWalletExtensionContainerFromConfig(config, logger)
	defer func() {
		err := walletExtContainer.Start()
		if err != nil {
			fmt.Printf("error stopping WE - %s", err)
		}
	}()

	go func() {
		err := walletExtContainer.Start()
		if err != nil {
			fmt.Printf("error in WE - %s", err)
		}
	}()

	walletExtensionAddr := fmt.Sprintf("%s:%d", common.Localhost, config.WalletExtensionPortHTTP)
	fmt.Printf("💡 Wallet extension started \n") // Some tests rely on seeing this message. Removed in next PR.
	fmt.Printf("💡 Obscuro Gateway started - visit http://%s to use it.\n", walletExtensionAddr)

	select {}
}
