package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

const (
	// Flag names, defaults and usages.
	nodeIDName    = "nodeID"
	nodeIDDefault = ""
	nodeIDUsage   = "The 20 bytes of the node's address (default \"\")"

	genesisName    = "isGenesis"
	genesisDefault = true
	genesisUsage   = "Whether the node is the first node to join the network (default true)"

	gossipRoundNanosName    = "gossipRoundNanos"
	gossipRoundNanosDefault = 8333
	gossipRoundNanosUsage   = "The duration of the gossip round (default 8333)"

	rpcTimeoutSecsName    = "rpcTimeoutSecs"
	rpcTimeoutSecsDefault = 3
	rpcTimeoutSecsUsage   = "The timeout for host <-> enclave RPC communication (default 3)"

	enclavePortName    = "enclavePort"
	enclavePortDefault = 10000
	enclavePortUsage   = "The port to use to connect to the Obscuro enclave service (default 10000)"

	helpCmd = "help"
	usage   = `CLI application for the ◠.bscuro host.

Usage:

    <executable> [flags]

The flags are:

  -%s string
    	%s
  -%s bool
    	%s
  -%s uint
    	%s
  -%s uint
    	%s
  -%s uint
    	%s`
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == helpCmd {
		usageFmt := fmt.Sprintf(usage, nodeIDName, nodeIDUsage, genesisName, genesisUsage,
			gossipRoundNanosName, gossipRoundNanosUsage, rpcTimeoutSecsName, rpcTimeoutSecsUsage, enclavePortName,
			enclavePortUsage)
		fmt.Println(usageFmt)
		return
	}

	nodeAddressBytes, isGenesis, gossipRoundNanos, rpcTimeoutSecs, enclavePort := parseCLIArgs()

	nodeID := common.BytesToAddress([]byte(*nodeAddressBytes))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: *gossipRoundNanos, ClientRPCTimeoutSecs: *rpcTimeoutSecs}
	enclaveClient := host.NewEnclaveRPCClient(*enclavePort, host.ClientRPCTimeoutSecs*time.Second)
	// TODO - Provide flags for our address and peer addresses
	aggP2P := p2p.NewP2P("localhost:10000", []string{})
	agg := host.NewAgg(nodeID, hostCfg, l1NodeDummy{}, nil, *isGenesis, enclaveClient, aggP2P)

	waitForEnclave(agg, *enclavePort)
	fmt.Printf("Connected to enclave server on port %d.\n", *enclavePort)
	agg.Start()
}

// Parses the CLI flags and arguments.
func parseCLIArgs() (*string, *bool, *uint64, *uint64, *uint64) {
	nodeID := flag.String(nodeIDName, nodeIDDefault, nodeIDUsage)
	genesis := flag.Bool(genesisName, genesisDefault, genesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(gossipRoundNanosDefault), gossipRoundNanosUsage)
	rpcTimeoutSecs := flag.Uint64(rpcTimeoutSecsName, rpcTimeoutSecsDefault, rpcTimeoutSecsUsage)
	enclavePort := flag.Uint64(enclavePortName, enclavePortDefault, enclavePortUsage)
	flag.Parse()

	return nodeID, genesis, gossipRoundNanos, rpcTimeoutSecs, enclavePort
}

// Waits for the enclave server to start, printing a wait message every two seconds.
func waitForEnclave(agg host.Node, enclavePort uint64) {
	i := 0
	for {
		if agg.Enclave.IsReady() == nil {
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
