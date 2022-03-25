package main

import (
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

func main() {
	setLogs()
	config := parseCLIArgs()

	nodeID := common.BytesToAddress([]byte(*config.nodeID))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: *config.gossipRoundNanos, ClientRPCTimeoutSecs: *config.rpcTimeoutSecs}
	enclaveClient := host.NewEnclaveRPCClient(*config.enclaveAddr, host.ClientRPCTimeoutSecs*time.Second)
	aggP2P := p2p.NewP2P(*config.ourP2PAddr, config.peerP2PAddrs)
	agg := host.NewAgg(nodeID, hostCfg, l1NodeDummy{}, nil, *config.isGenesis, enclaveClient, aggP2P)

	waitForEnclave(agg, *config.enclaveAddr)
	agg.Start()
}

// Sets the log file.
func setLogs() {
	logFile, err := os.Create("host_logs.txt")
	if err != nil {
		panic(err)
	}
	log.SetLog(logFile)
}

// Waits for the enclave server to start, printing a wait message every two seconds.
func waitForEnclave(agg host.Node, enclaveAddr string) {
	i := 0
	for {
		if agg.Enclave.IsReady() == nil {
			fmt.Printf("Connected to enclave server on address %s.\n", enclaveAddr)
			break
		}
		time.Sleep(100 * time.Millisecond)
		i++

		if i >= 20 {
			fmt.Printf("Trying to connect to enclave server on address %s...\n", enclaveAddr)
			i = 0
		}
	}
}

// TODO - Replace this dummy once we have implemented communication with L1 nodes.
type l1NodeDummy struct{}

func (l l1NodeDummy) RPCBlockchainFeed() []*types.Block {
	return []*types.Block{}
}

func (l l1NodeDummy) BroadcastTx(obscurocommon.EncodedL1Tx) {}
