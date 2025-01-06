package container

import (
	"context"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave"

	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	obscuroGenesis "github.com/ten-protocol/go-ten/go/enclave/genesis"
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
	logger := log.New(log.EnclaveCmp, config.LogLevel, config.LogPath, log.NodeIDKey, config.NodeID)

	logger.Info(fmt.Sprintf("Building enclave container with config: %+v", config))

	return NewEnclaveContainerWithLogger(config, logger)
}

// NewEnclaveContainerWithLogger is useful for testing etc.
func NewEnclaveContainerWithLogger(config *enclaveconfig.EnclaveConfig, logger gethlog.Logger) *EnclaveContainer {
	contractAddr := config.ManagementContractAddress
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&contractAddr, logger)

	genesis, err := obscuroGenesis.New(config.TenGenesis)
	if err != nil {
		logger.Crit("unable to parse obscuro genesis", log.ErrKey, err)
	}

	encl := enclave.NewEnclave(config, genesis, mgmtContractLib, logger)
	rpcServer := enclave.NewEnclaveRPCServer(config.RPCAddress, encl, logger)

	return &EnclaveContainer{
		Enclave:   encl,
		RPCServer: rpcServer,
		Logger:    logger,
	}
}
