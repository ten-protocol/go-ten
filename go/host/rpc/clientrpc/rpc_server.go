package clientrpc

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	allOrigins = "*"
)

// Server is the layer responsible for handling RPC requests from Obscuro client applications.
type Server interface {
	Start()
	Stop()
	RegisterAPIs(apis []rpc.API)
}

// An implementation of `host.Server` that reuses the Geth `node` package for client communication.
type serverImpl struct {
	node   *node.Node
	logger gethlog.Logger
}

func NewServer(config *config.HostConfig, logger gethlog.Logger) Server {
	rpcConfig := node.Config{
		Logger: logger.New(log.CmpKey, log.HostRPCCmp),
	}
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
		logger.Crit("could not create new client server.", log.ErrKey, err)
	}

	return &serverImpl{node: rpcServerNode, logger: logger}
}

func (s *serverImpl) RegisterAPIs(apis []rpc.API) {
	s.node.RegisterAPIs(apis)
}

func (s *serverImpl) Start() {
	// make sure you *have* registered apis
	if err := s.node.Start(); err != nil {
		s.logger.Crit("could not start node client server.", log.ErrKey, err)
	}
}

func (s *serverImpl) Stop() {
	err := s.node.Close()
	if err != nil {
		s.logger.Crit("could not stop node client server.", log.ErrKey, err)
	}
}
