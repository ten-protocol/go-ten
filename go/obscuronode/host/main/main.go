package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"os"
	"time"
)

const (
	nodeAddressFlag       = "nodeAddress"
	nodeAddressUsage      = "The 20 bytes of the node's address"
	genesisFlag           = "isGenesis"
	genesisUsage          = "Whether the node is the first node to join the network"
	gossipRoundNanosFlag  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round"
	rpcTimeoutSecsFlag    = "rpcTimeoutSecs"
	rpcTimeoutSecsUsage   = "The timeout for host <-> enclave RPC communication"
	enclavePortFlag       = "enclavePort"
	enclavePortUsage      = "The port to use to connect to the Obscuro enclave service"
	usage                 = `CLI application for the ◠.bscuro host.

Usage:

    <executable> [flags]

The flags are:

    -%s   string   %s
    -%s   bool   %s
    -%s   int   %s
    -%s   int   %s
    -%s   int   %s`
)

func main() {
	if len(os.Args) == 1 {
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
	agg := host.NewAgg(nodeAddress, hostCfg, l1NodeDummy{}, &l2Network, nil, *isGenesis, *enclavePort)

	waitForEnclave(agg, *enclavePort)
	fmt.Printf("Connected to enclave server on port %d.\n", *enclavePort)
	agg.Start()
}

// Parses the CLI flags and arguments.
func parseCLIArgs() (*string, *bool, *uint64, *uint64, *uint64) {
	var nodeAddressBytes = flag.String(nodeAddressFlag, "", nodeAddressUsage)
	var genesis = flag.Bool(genesisFlag, true, genesisUsage)
	var gossipRoundNanos = flag.Uint64(gossipRoundNanosFlag, uint64(25_000/3), gossipRoundNanosUsage)
	var rpcTimeoutSecs = flag.Uint64(rpcTimeoutSecsFlag, 3, rpcTimeoutSecsUsage)
	var enclavePort = flag.Uint64(enclavePortFlag, 10000, enclavePortUsage)
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

// TODO - Replace this dummy with actual node interaction once we have implemented P2P and gossiping.
type l2NetworkDummy struct{}

func (l *l2NetworkDummy) BroadcastRollup(obscurocommon.EncodedRollup) {}
func (l *l2NetworkDummy) BroadcastTx(nodecommon.EncryptedTx)          {}

// todo - joel - explain why
type l1NodeDummy struct{}

func (l l1NodeDummy) RPCBlockchainFeed() []*types.Block {
	return []*types.Block{}
}

func (l l1NodeDummy) BroadcastTx(obscurocommon.EncodedL1Tx) {}
