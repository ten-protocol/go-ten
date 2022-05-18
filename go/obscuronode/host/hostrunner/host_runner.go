package hostrunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
)

// RunHost runs an Obscuro host as a standalone process.
func RunHost(config HostConfig) {
	setLogs(config.LogPath)

	nodeID := common.BytesToAddress([]byte(config.NodeID))
	hostCfg := host.AggregatorCfg{
		GossipRoundDuration:  time.Duration(config.GossipRoundNanos),
		ClientRPCTimeoutSecs: config.RPCTimeoutSecs,
		HasRPC:               true,
		RPCAddress:           &config.ClientServerAddr,
	}

	nodeWallet := wallet.NewInMemoryWallet(config.PrivateKeyString)
	contractAddr := common.HexToAddress(config.ContractAddress)
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&contractAddr)
	erc20ContracLib := erc20contractlib.NewERC20ContractLib(&contractAddr)
	ethWallet := datagenerator.RandomWallet()

	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethclient.NewEthClient(nodeID, config.EthClientHost, uint(config.EthClientPort), nodeWallet, &contractAddr)
	if err != nil {
		panic(err)
	}

	enclaveClient := host.NewEnclaveRPCClient(config.EnclaveAddr, host.ClientRPCTimeoutSecs*time.Second, nodeID)
	aggP2P := p2p.NewSocketP2PLayer(config.OurP2PAddr, config.PeerP2PAddrs, nodeID)
	agg := host.NewObscuroAggregator(nodeID, hostCfg, nil, config.IsGenesis, aggP2P, l1Client, enclaveClient, ethWallet, mgmtContractLib, erc20ContracLib)

	fmt.Println("Starting Obscuro host...")
	agg.Start()

	handleInterrupt(agg)
}

// Sets the log file.
func setLogs(logPath string) {
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	log.SetLog(logFile)
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
