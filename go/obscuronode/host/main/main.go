package main

import (
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
	enclaveClient := host.NewEnclaveRPCClient(*config.enclaveAddr, host.ClientRPCTimeoutSecs*time.Second, nodeID)
	aggP2P := p2p.NewSocketP2PLayer(*config.ourP2PAddr, config.peerP2PAddrs)
	agg := host.NewObscuroAggregator(nodeID, hostCfg, l1NodeDummy{}, nil, *config.isGenesis, enclaveClient, aggP2P)

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

// TODO - Replace this dummy once we have implemented communication with L1 nodes.
type l1NodeDummy struct{}

func (l l1NodeDummy) RPCBlockchainFeed() []*types.Block {
	return []*types.Block{}
}

func (l l1NodeDummy) BroadcastTx(obscurocommon.EncodedL1Tx) {}
