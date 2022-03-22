package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"os"
	"strconv"
)

const (
	nodeAddressFlag  = "nodeAddress"
	nodeAddressUsage = "The 20 bytes of the node's address"
	usage            = `CLI application for the ◠.bscuro host. Usage: <executable flag1 ... flagN arg1 ... argN>

Flags:
-%s   string   %s

Arguments:
  port   The port on which to serve the Obscuro enclave service`
)

func main() {
	nodeAddressBytes, port, isInvalid := parseCLIArgs()
	if isInvalid {
		usageFmt := fmt.Sprintf(usage, nodeAddressFlag, nodeAddressUsage)
		fmt.Println(usageFmt)
		return
	}

	nodeAddress := common.BytesToAddress([]byte(*nodeAddressBytes))
	enclaveServer, err := enclave.StartServer(port, nodeAddress, nil) // todo - handle error
	println("jjj")
	println(err)
	for {
		
	}
	defer enclaveServer.Stop()
}

// Parses the CLI flags and arguments.
func parseCLIArgs() (*string, uint64, bool) {
	var nodeAddressBytes = flag.String(nodeAddressFlag, "", nodeAddressUsage)
	flag.Parse()

	if flag.NArg() != 1 {
		return nil, 0, true
	}

	enclavePort, err := strconv.ParseUint(os.Args[len(os.Args)-1], 10, 64)
	if err != nil {
		return nil, 0, true
	}

	return nodeAddressBytes, enclavePort, false
}
