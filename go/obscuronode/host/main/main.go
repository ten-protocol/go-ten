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
	enclavePortUsage      = "The address to use to connect to the Obscuro enclave service (default \"localhost:10000\")"
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

	nodeAddressBytes, isGenesis, gossipRoundNanos, rpcTimeoutSecs, enclaveAddress := parseCLIArgs()

	nodeAddress := common.BytesToAddress([]byte(*nodeAddressBytes))
	hostCfg := host.AggregatorCfg{GossipRoundDuration: *gossipRoundNanos, ClientRPCTimeoutSecs: *rpcTimeoutSecs}
	enclaveClient := host.NewEnclaveRPCClient(*enclaveAddress, host.ClientRPCTimeoutSecs*time.Second)
	// TODO - Provide flags for our address and peer addresses
	aggP2P := p2p.NewP2P("localhost:10000", []string{})
	agg := host.NewAgg(nodeAddress, hostCfg, l1NodeDummy{}, nil, *isGenesis, enclaveClient, aggP2P)

	waitForEnclave(agg, *enclaveAddress)
	fmt.Printf("Connected to enclave server on port %s.\n", *enclaveAddress)
	agg.Start()
}

// Parses the CLI flags and arguments.
func parseCLIArgs() (*string, *bool, *uint64, *uint64, *string) {
	nodeAddressBytes := flag.String(nodeAddressFlag, "", nodeAddressUsage)
	genesis := flag.Bool(genesisFlag, true, genesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosFlag, uint64(8333), gossipRoundNanosUsage)
	rpcTimeoutSecs := flag.Uint64(rpcTimeoutSecsFlag, 3, rpcTimeoutSecsUsage)
	enclaveAddress := flag.String(enclavePortFlag, "localhost:10000", enclavePortUsage)
	flag.Parse()

	return nodeAddressBytes, genesis, gossipRoundNanos, rpcTimeoutSecs, enclaveAddress
}

// Waits for the enclave server to start, printing a wait message every two seconds.
func waitForEnclave(agg host.Node, enclaveAddress string) {
	i := 0
	for {
		if agg.Enclave.IsReady() == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
		i++

		if i >= 20 {
			fmt.Printf("Trying to connect to enclave on address %s...\n", enclaveAddress)
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
