package main

import (
	"fmt"
	stats2 "github.com/obscuronet/obscuro-playground/integration/simulation/stats"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

func main() {
	config := parseCLIArgs()

	// todo - joel - override ethClientHost as 20.26.233.214, ethClientPort as 12000, and clientServerAddress as 20.26.233.214:13000
	// todo - joel - need to expose eth port (internally on 12000) to inject txs

	nodeID := common.BytesToAddress([]byte(config.nodeID))
	nodeWallet := wallet.NewInMemoryWallet(config.privateKeyString)
	contractAddr := common.HexToAddress(config.contractAddress)
	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethclient.NewEthClient(nodeID, config.ethClientHost, uint(config.ethClientPort), nodeWallet, contractAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to Obscuro host...")
	l2Client := obscuroclient.NewClient(config.clientServerAddr)

	// todo - joel - parameterise number of nodes
	stats := stats2.NewStats(2)
	txInjector := simulation.NewTransactionInjector(
		3,
		time.Second,
		stats,
		[]ethclient.EthClient{l1Client},
		[]*obscuroclient.Client{&l2Client},
	)

	fmt.Println("Injecting transactions...")
	txInjector.Start()
}
