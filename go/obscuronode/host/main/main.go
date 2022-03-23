package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const (
	helpCmd               = "help"
	nodeAddressFlag       = "nodeAddress"
	nodeAddressUsage      = "The 20 bytes of the node's address (default \"\")"
	genesisFlag           = "isGenesis"
	genesisUsage          = "Whether the node is the first node to join the network (default true)"
	gossipRoundNanosFlag  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round (default 8333)"
	rpcTimeoutSecsFlag    = "rpcTimeoutSecs"
	rpcTimeoutSecsUsage   = "The timeout for host <-> enclave RPC communication (default 3)"
	enclavePortFlag       = "enclavePort"
	enclavePortUsage      = "The port to use to connect to the Obscuro enclave service (default 10000)"
	usage                 = `CLI application for the ◠.bscuro host.

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
		usageFmt := fmt.Sprintf(usage, nodeAddressFlag, nodeAddressUsage, genesisFlag, genesisUsage,
			gossipRoundNanosFlag, gossipRoundNanosUsage, rpcTimeoutSecsFlag, rpcTimeoutSecsUsage, enclavePortFlag,
			enclavePortUsage)
		fmt.Println(usageFmt)
		return
	}

	nodeAddressBytes, isGenesis, gossipRoundNanos, rpcTimeoutSecs, enclavePort := parseCLIArgs()

	nodeAddress := common.BytesToAddress([]byte(*nodeAddressBytes))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: *gossipRoundNanos, ClientRPCTimeoutSecs: *rpcTimeoutSecs}
	l2Network := l2NetworkDummy{}
	enclaveClient := host.NewEnclaveRPCClient(*enclavePort, host.ClientRPCTimeoutSecs*time.Second)
	// todo - joel - use flag for p2p address
	agg := host.NewAgg(nodeAddress, hostCfg, l1NodeDummy{}, &l2Network, nil, *isGenesis, enclaveClient, "localhost:10000")

	waitForEnclave(agg, *enclavePort)
	fmt.Printf("Connected to enclave server on port %d.\n", *enclavePort)
	agg.Start()
}

// Parses the CLI flags and arguments.
func parseCLIArgs() (*string, *bool, *uint64, *uint64, *uint64) {
	nodeAddressBytes := flag.String(nodeAddressFlag, "", nodeAddressUsage)
	genesis := flag.Bool(genesisFlag, true, genesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosFlag, uint64(8333), gossipRoundNanosUsage)
	rpcTimeoutSecs := flag.Uint64(rpcTimeoutSecsFlag, 3, rpcTimeoutSecsUsage)
	enclavePort := flag.Uint64(enclavePortFlag, 10000, enclavePortUsage)
	flag.Parse()

	return nodeAddressBytes, genesis, gossipRoundNanos, rpcTimeoutSecs, enclavePort
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

// TODO - Replace this dummy once we have implemented P2P communication and gossiping between L2 nodes.
type l2NetworkDummy struct{}

func (l *l2NetworkDummy) BroadcastRollup(obscurocommon.EncodedRollup) {}
func (l *l2NetworkDummy) BroadcastTx(nodecommon.EncryptedTx)          {}

// TODO - Replace this dummy once we have implemented communication with L1 nodes.
type l1NodeDummy struct{}

func (l l1NodeDummy) RPCBlockchainFeed() []*types.Block {
	return []*types.Block{}
}

func (l l1NodeDummy) BroadcastTx(obscurocommon.EncodedL1Tx) {}
