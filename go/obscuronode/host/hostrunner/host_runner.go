package hostrunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"
)

const ClientRPCTimeoutSecs = 5

// RunHost runs an Obscuro host as a standalone process.
func RunHost(config HostConfig) {
	contractAddr := common.HexToAddress(config.ContractAddress)

	// todo - joel - rebuild this around the host.Config object
	hostConfig := host.Config{
		ID:                    common.BytesToAddress([]byte(config.NodeID)),
		IsGenesis:             config.IsGenesis,
		GossipRoundDuration:   time.Duration(config.GossipRoundNanos),
		ClientRPCTimeoutSecs:  config.RPCTimeoutSecs,
		HasClientRPC:          true,
		ClientRPCAddress:      &config.ClientServerAddr,
		EnclaveRPCAddress:     &config.EnclaveAddr,
		EnclaveRPCTimeout:     ClientRPCTimeoutSecs * time.Second, // todo - joel - pass this via CLI
		P2PAddress:            &config.OurP2PAddr,
		AllP2PAddresses:       config.PeerP2PAddrs,
		L1NodeHost:            &config.EthClientHost,
		L1NodeWebsocketPort:   uint(config.EthClientPort),
		RollupContractAddress: &contractAddr,
	}

	nodeWallet := wallet.NewInMemoryWallet(config.PrivateKeyString)
	txHandler := mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr)

	fmt.Println("Connecting to L1 network...")
	log.Info("Connecting to L1 network...")
	l1Client, err := host.NewEthClient(hostConfig, nodeWallet)
	if err != nil {
		log.Panic("could not create Ethereum client. Cause: %s", err)
	}

	enclaveClient := host.NewEnclaveRPCClient(hostConfig)
	aggP2P := p2p.NewSocketP2PLayer(hostConfig)
	agg := host.NewHost(hostConfig, nil, aggP2P, l1Client, enclaveClient, txHandler)

	fmt.Println("Starting Obscuro host...")
	log.Info("Starting Obscuro host...")
	agg.Start()

	handleInterrupt(agg)
}

// SetLogs sets the log file.
func SetLogs(logPath string) {
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(fmt.Sprintf("could not create log file. Cause: %s", err))
	}
	log.OutputToFile(logFile)
}

// Shuts down the Obscuro host when an interrupt is received.
func handleInterrupt(host host.Node) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	host.Stop()
	fmt.Println("Obscuro host stopping...")
	os.Exit(1)
}
