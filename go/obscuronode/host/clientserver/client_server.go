package clientserver

import (
	"fmt"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"strconv"
	"strings"
)

const apiNamespaceEth = "eth"
const apiVersion1 = "1.0"

// An implementation of `host.ClientServer` that reuses the Geth `node` package for client communication.
type clientServerImpl struct {
	node *node.Node
}

// NewClientServer returns a `host.ClientServer` that wraps the Geth `node` package for client communication, and
// offers `NewEthAPI` under the "eth" namespace.
func NewClientServer(address string) host.ClientServer {
	hostAndPort := strings.Split(address, ":")
	if len(hostAndPort) != 2 {
		panic(fmt.Sprintf("Client server expected address in the form <host>:<port>, but received %s", address))
	}
	port, err := strconv.Atoi(hostAndPort[1])
	if err != nil {
		panic(fmt.Sprintf("Client server port %s could not be converted to an integer", hostAndPort[1]))
	}

	nodeConfig := node.Config{
		// We do not listen over websockets and IPC for now.
		HTTPHost: hostAndPort[0],
		HTTPPort: port,
	}
	clientServerNode, err := node.New(&nodeConfig)
	if err != nil {
		panic(err)
	}

	rpcAPIs := []rpc.API{
		{
			Namespace: apiNamespaceEth,
			Version:   apiVersion1,
			Service:   NewEthAPI(),
			Public:    true,
		},
	}
	clientServerNode.RegisterAPIs(rpcAPIs)

	return clientServerImpl{
		node: clientServerNode,
	}
}

func (server clientServerImpl) Start() {
	err := server.node.Start()
	if err != nil {
		panic(err)
	}
}

func (server clientServerImpl) Stop() {
	err := server.node.Close()
	if err != nil {
		return // todo - joel - more graceful approach?
	}
}
