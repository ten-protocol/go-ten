package main

import (
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

func main() {
	config := parseCLIArgs()

	nodeID := common.BytesToAddress([]byte(*config.nodeID))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: *config.gossipRoundNanos, ClientRPCTimeoutSecs: *config.rpcTimeoutSecs}
	enclaveClient := host.NewEnclaveRPCClient(*config.enclavePort, host.ClientRPCTimeoutSecs*time.Second)
	aggP2P := p2p.NewP2P(*config.ourP2PAddr, config.peerP2PAddrs)
	agg := host.NewAgg(nodeID, hostCfg, l1NodeDummy{}, nil, *config.isGenesis, enclaveClient, aggP2P)

	waitForEnclave(agg, *config.enclavePort)
	agg.Start()
}

// Waits for the enclave server to start, printing a wait message every two seconds.
func waitForEnclave(agg host.Node, enclavePort uint64) {
	i := 0
	for {
		if agg.Enclave.IsReady() == nil {
			fmt.Printf("Connected to enclave server on port %d.\n", enclavePort)
			break
		}
		time.Sleep(100 * time.Millisecond)
		i++

		if i >= 20 {
			fmt.Printf("Trying to connect to enclave on port %d...\n", enclavePort)
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
