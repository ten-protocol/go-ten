package clientrpc

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	allOrigins = "*"

	APIVersion1             = "1.0"
	APINamespaceObscuro     = "obscuro"
	APINamespaceEth         = "eth"
	APINamespaceObscuroScan = "obscuroscan"
	APINamespaceScan        = "scan"
	APINamespaceNetwork     = "net"
	APINamespaceTest        = "test"
	APINamespaceDebug       = "debug"
)

// Server is the layer responsible for handling RPC requests from Obscuro client applications.
// - reuses the Geth `node` package for client communication
// - implements the host.Service interface
type Server struct {
	node   *node.Node
	logger gethlog.Logger
}

// CreateServerFactory returns a factory that creates a new RPC server for the host, configuring APIs based on the provided config
// - the factory is created with a hostStop function so that RPC APIs can stop the host
// todo (@matt) change the tests so that this RPC method is not needed as it's very unintuitive
func CreateServerFactory(hostStop func() error) host.ServiceFactory[host.RPCServerService] {
	return func(config *config.HostConfig, serviceLocator host.ServiceLocator, logger gethlog.Logger) (host.RPCServerService, error) {
		if !(config.HasClientRPCHTTP || config.HasClientRPCWebsockets) {
			return nil, fmt.Errorf("cannot create RPC server as RPC is not enabled in host config")
		}
		apis := []rpc.API{
			{
				Namespace: APINamespaceObscuro,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroAPI(serviceLocator),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewEthereumAPI(config.ObscuroChainID, serviceLocator, logger.New(log.CmpKey, "eth-rpc")),
				Public:    true,
			},
			{
				Namespace: APINamespaceObscuroScan,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroScanAPI(serviceLocator),
				Public:    true,
			},
			{
				Namespace: APINamespaceNetwork,
				Version:   APIVersion1,
				Service:   clientapi.NewNetworkAPI(config.ObscuroChainID),
				Public:    true,
			},
			{
				Namespace: APINamespaceTest,
				Version:   APIVersion1,
				Service:   clientapi.NewTestAPI(hostStop),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewFilterAPI(serviceLocator, logger.New(log.CmpKey, "filter-rpc")),
				Public:    true,
			},
			{
				Namespace: APINamespaceScan,
				Version:   APIVersion1,
				Service:   clientapi.NewScanAPI(serviceLocator, logger.New(log.CmpKey, "scan-rpc")),
				Public:    true,
			},
		}
		if config.DebugNamespaceEnabled {
			apis = append(apis, rpc.API{
				Namespace: APINamespaceDebug,
				Version:   APIVersion1,
				Service:   clientapi.NewNetworkDebug(serviceLocator),
				Public:    true,
			})
		}

		return NewServer(config, apis, logger)
	}
}

func NewServer(config *config.HostConfig, apis []rpc.API, logger gethlog.Logger) (*Server, error) {
	rpcConfig := node.Config{
		Logger: logger.New(log.CmpKey, log.HostRPCCmp),
	}
	if config.HasClientRPCHTTP {
		rpcConfig.HTTPHost = config.ClientRPCHost
		rpcConfig.HTTPPort = int(config.ClientRPCPortHTTP)
		// todo (@pedro) - review if this poses a security issue
		rpcConfig.HTTPVirtualHosts = []string{allOrigins}
	}
	if config.HasClientRPCWebsockets {
		rpcConfig.WSHost = config.ClientRPCHost
		rpcConfig.WSPort = int(config.ClientRPCPortWS)
		// todo (@pedro) - review if this poses a security issue
		rpcConfig.WSOrigins = []string{allOrigins}
	}

	rpcServerNode, err := node.New(&rpcConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create new client server: %w", err)
	}

	rpcServerNode.RegisterAPIs(apis)

	return &Server{node: rpcServerNode, logger: logger}, nil
}

func (s *Server) Start() error {
	return s.node.Start()
}

func (s *Server) Stop() {
	err := s.node.Close()
	if err != nil {
		s.logger.Error("error closing RPC server", log.ErrKey, err)
	}
}

func (s *Server) HealthStatus() hostcommon.HealthStatus {
	// always return healthy for now, this method means we satisfy the host Service interface
	return &hostcommon.BasicErrHealthStatus{ErrMsg: ""}
}
