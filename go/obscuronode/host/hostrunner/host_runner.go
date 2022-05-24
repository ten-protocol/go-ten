package hostrunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"
)

const ClientRPCTimeoutSecs = 5

// RunHost runs an Obscuro host as a standalone process.
func RunHost(config HostConfig) {
	nodeID := common.BytesToAddress([]byte(config.NodeID))
	// todo - joel - rebuild this around the host.Config object
	hostCfg := host.Config{
		ID:                   nodeID,
		IsGenesis:            config.IsGenesis,
		GossipRoundDuration:  time.Duration(config.GossipRoundNanos),
		ClientRPCTimeoutSecs: config.RPCTimeoutSecs,
		HasClientRPC:         true,
		ClientRPCAddress:     &config.ClientServerAddr,
		EnclaveRPCAddress:    &config.EnclaveAddr,
		EnclaveRPCTimeout:    ClientRPCTimeoutSecs * time.Second, // todo - joel - pass this via CLI
	}

	nodeWallet := wallet.NewInMemoryWallet(config.PrivateKeyString)
	contractAddr := common.HexToAddress(config.ContractAddress)
	txHandler := mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr)

	fmt.Println("Connecting to L1 network...")
	log.Info("Connecting to L1 network...")
	l1Client, err := ethclient.NewEthClient(nodeID, config.EthClientHost, uint(config.EthClientPort), nodeWallet, contractAddr)
	if err != nil {
		log.Panic("could not create Ethereum client. Cause: %s", err)
	}

	enclaveClient := host.NewEnclaveRPCClient(hostCfg)
	aggP2P := p2p.NewSocketP2PLayer(config.OurP2PAddr, config.PeerP2PAddrs, nodeID)
	agg := host.NewHost(hostCfg, nil, aggP2P, l1Client, enclaveClient, txHandler)

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
