package enclavecontainer

import (
	"context"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/container"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

// TODO - Replace with the genesis.json of Obscuro's L1 network.
const hardcodedGenesisJSON = "TODO - REPLACE ME"

type EnclaveContainer struct {
	Enclave   common.Enclave
	RpcServer *enclave.EnclaveRpcServer
	Logger    gethlog.Logger
}

// NewEnclaveContainer wires up the components of the Enclave and its RPC server. Manages their lifecycle/monitors their status
func NewEnclaveContainer(config config.EnclaveConfig) *EnclaveContainer {
	// todo - improve this wiring, perhaps setup DB etc. at this level and inject into enclave
	logger := log.New(log.EnclaveCmp, config.LogLevel, config.LogPath, log.NodeIDKey, config.HostID)

	contractAddr := config.ManagementContractAddress
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&contractAddr, logger)
	erc20ContractLib := erc20contractlib.NewERC20ContractLib(&contractAddr, config.ERC20ContractAddresses...)

	if config.ValidateL1Blocks {
		config.GenesisJSON = []byte(hardcodedGenesisJSON)
	}
	encl := enclave.NewEnclave(config, mgmtContractLib, erc20ContractLib, logger)
	rpcServer := enclave.NewEnclaveRPCServer(config.Address, encl, logger)

	return &EnclaveContainer{
		Enclave:   encl,
		RpcServer: rpcServer,
		Logger:    logger,
	}
}

func (e *EnclaveContainer) Start() error {
	err := e.RpcServer.StartServer()
	if err != nil {
		return err
	}
	e.Logger.Info("Obscuro enclave service started.")
	return nil
}

func (e *EnclaveContainer) Stop() error {
	_, err := e.RpcServer.Stop(context.Background(), nil)
	if err != nil {
		e.Logger.Warn("unable to cleanly stop enclave", log.ErrKey, err)
		return err
	}
	return nil
}

func (e *EnclaveContainer) Status() container.Status {
	// todo: return recovery if DB is unavailable or if host is not keeping up with its responsibilities
	return container.Running
}
