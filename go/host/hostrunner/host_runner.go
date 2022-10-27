package hostrunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/go-obscuro/go/host/node"

	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
)

// RunHost runs an Obscuro host as a standalone process.
func RunHost(config config.HostConfig) {
	logger := log.New(log.HostCmp, config.LogLevel, config.LogPath, log.NodeIDKey, config.ID)

	logger.Info(fmt.Sprintf("Starting node with config: %+v", config))
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&config.RollupContractAddress, logger)

	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethadapter.NewEthClient(config.L1NodeHost, config.L1NodeWebsocketPort, config.L1RPCTimeout, config.ID, logger)
	if err != nil {
		logger.Crit("could not create Ethereum client.", log.ErrKey, err)
	}

	ethWallet := wallet.NewInMemoryWalletFromConfig(config, logger)
	nonce, err := l1Client.Nonce(ethWallet.Address())
	if err != nil {
		logger.Crit("could not retrieve Ethereum account nonce.", log.ErrKey, err)
	}
	ethWallet.SetNonce(nonce)

	enclaveClient := enclaverpc.NewClient(config, logger)
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	aggP2P := p2p.NewSocketP2PLayer(config, p2pLogger)
	agg := node.NewHost(config, nil, aggP2P, l1Client, enclaveClient, ethWallet, mgmtContractLib, logger)

	fmt.Println("Starting Obscuro host...")
	logger.Info("Starting Obscuro host...")
	agg.Start()

	handleInterrupt(agg)
}

// Shuts down the Obscuro host when an interrupt is received.
func handleInterrupt(host host.Host) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	host.Stop()
	fmt.Println("Obscuro host stopping...")
	os.Exit(1)
}
