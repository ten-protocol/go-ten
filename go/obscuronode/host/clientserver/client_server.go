package clientserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

const (
	apiNamespaceEth = "eth"
	apiVersion1     = "1.0"
)

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
		panic(fmt.Sprintf("Could not create new client server. Cause: %s", err))
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
	if err := server.node.Start(); err != nil {
		panic(fmt.Sprintf("Could not start node client server. Cause: %s", err))
	}
}

func (server clientServerImpl) Stop() {
	if err := server.node.Close(); err != nil {
		return
	}
}
