package container

import (
	"context"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave"

	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
)

type EnclaveContainer struct {
	Enclave   common.Enclave
	RPCServer *enclave.RPCServer
	Logger    gethlog.Logger
}

func (e *EnclaveContainer) Start() error {
	err := e.RPCServer.StartServer()
	if err != nil {
		return err
	}
	e.Logger.Info("obscuro enclave RPC service started.")
	return nil
}

func (e *EnclaveContainer) Stop() error {
	_, err := e.RPCServer.Stop(context.Background(), nil)
	if err != nil {
		e.Logger.Error("Unable to cleanly stop enclave", log.ErrKey, err)
		return err
	}
	return nil
}

// NewEnclaveContainerFromConfig wires up the components of the Enclave and its RPC server. Manages their lifecycle/monitors their status
func NewEnclaveContainerFromConfig(config *enclaveconfig.EnclaveConfig) *EnclaveContainer {
	// todo - improve this wiring, perhaps setup DB etc. at this level and inject into enclave
	// (at that point the WithLogger constructor could be a full DI constructor like the HostContainer tries, for testability)
	logger := log.New(log.EnclaveCmp, config.LogLevel, config.LogPath, log.NodeIDKey, config.HostID)

	// todo - this is for debugging purposes only, should be remove in the future
	fmt.Printf("Building enclave container with config: %+v\n", config)
	logger.Info(fmt.Sprintf("Building enclave container with config: %+v", config))

	return NewEnclaveContainerWithLogger(config, logger)
}

// NewEnclaveContainerWithLogger is useful for testing etc.
func NewEnclaveContainerWithLogger(config *enclaveconfig.EnclaveConfig, logger gethlog.Logger) *EnclaveContainer {
	encl := enclave.NewEnclave(config, logger)
	rpcServer := enclave.NewEnclaveRPCServer(config.Address, encl, logger)

	return &EnclaveContainer{
		Enclave:   encl,
		RPCServer: rpcServer,
		Logger:    logger,
	}
}
