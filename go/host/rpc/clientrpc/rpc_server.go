package clientrpc

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
)

const (
	allOrigins = "*"
)

// Server is the layer responsible for handling RPC requests from Obscuro client applications.
type Server interface {
	Start()
	Stop()
}

// An implementation of `host.Server` that reuses the Geth `node` package for client communication.
type serverImpl struct {
	node *node.Node
}

func NewServer(config config.HostConfig, rpcAPIs []rpc.API) Server {
	rpcConfig := node.Config{}
	if config.HasClientRPCHTTP {
		rpcConfig.HTTPHost = config.ClientRPCHost
		rpcConfig.HTTPPort = int(config.ClientRPCPortHTTP)
		// TODO review if this poses a security issue
		rpcConfig.HTTPVirtualHosts = []string{allOrigins}
	}
	if config.HasClientRPCWebsockets {
		rpcConfig.WSHost = config.ClientRPCHost
		rpcConfig.WSPort = int(config.ClientRPCPortWS)
		// TODO review if this poses a security issue
		rpcConfig.WSOrigins = []string{allOrigins}
	}

	rpcServerNode, err := node.New(&rpcConfig)
	if err != nil {
		log.Panic("could not create new client server. Cause: %s", err)
	}
	rpcServerNode.RegisterAPIs(rpcAPIs)

	return &serverImpl{node: rpcServerNode}
}

func (s *serverImpl) Start() {
	if err := s.node.Start(); err != nil {
		log.Panic("could not start node client server. Cause: %s", err)
	}
}

func (s *serverImpl) Stop() {
	s.node.Server().Stop()
	go func() {
		// TODO - Investigate why this never returns.
		_ = s.node.Close()
	}()
}
