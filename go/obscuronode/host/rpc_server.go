package host

import (
	"strconv"
	"strings"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	apiNamespaceObscuro  = "obscuro"
	apiNamespaceEthereum = "eth"
	apiVersion1          = "1.0"
)

// An implementation of `host.RPCServer` that reuses the Geth `node` package for client communication.
type rpcServerImpl struct {
	node *node.Node
}

func NewRPCServer(address string, host *Node) RPCServer {
	hostAndPort := strings.Split(address, ":")
	if len(hostAndPort) != 2 {
		log.Panic("client server expected address in the form <host>:<port>, but received %s", address)
	}
	port, err := strconv.Atoi(hostAndPort[1])
	if err != nil {
		log.Panic("client server port %s could not be converted to an integer", hostAndPort[1])
	}

	nodeConfig := node.Config{
		// We do not listen over websockets and IPC for now.
		HTTPHost: hostAndPort[0],
		HTTPPort: port,
		WSHost:   hostAndPort[0],
		WSPort:   3102, // todo - joel - parameterise
		// todo - joel - add origins policy if needed
	}
	rpcServerNode, err := node.New(&nodeConfig)
	if err != nil {
		log.Panic("could not create new client server. Cause: %s", err)
	}

	rpcAPIs := []rpc.API{
		{
			Namespace: apiNamespaceObscuro,
			Version:   apiVersion1,
			Service:   NewObscuroAPI(host),
			Public:    true,
		},
		{
			Namespace: apiNamespaceEthereum,
			Version:   apiVersion1,
			Service:   NewEthereumAPI(),
			Public:    true,
		},
	}
	rpcServerNode.RegisterAPIs(rpcAPIs)

	return rpcServerImpl{node: rpcServerNode}
}

func (s rpcServerImpl) Start() {
	if err := s.node.Start(); err != nil {
		log.Panic("could not start node client server. Cause: %s", err)
	}
}

func (s rpcServerImpl) Stop() error {
	return s.node.Close()
}
