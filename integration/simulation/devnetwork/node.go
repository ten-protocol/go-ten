package devnetwork

import (
	"fmt"
	"math/big"
	"os"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ten-protocol/go-ten/go/host"

	"github.com/ten-protocol/go-ten/go/enclave/storage/init/sqlite"

	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/metrics"
	"github.com/ten-protocol/go-ten/go/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
	"github.com/ten-protocol/go-ten/go/host/p2p"
	"github.com/ten-protocol/go-ten/go/host/rpc/clientrpc"
	"github.com/ten-protocol/go-ten/go/host/rpc/enclaverpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

// InMemNodeOperator represents an Obscuro node playing a role in a DevSimulation
//
// Anything a node operator could do, that we need to simulate belongs here, for example:
// * start/stop/reset the host/enclave components of the node
// * provide convenient access to RPC clients for the node
// * it might even provide access to things like the database files for the node if needed
//
// Note: InMemNodeOperator will panic when things go wrong, we want to fail fast in sims and avoid verbose error handling in usage
type InMemNodeOperator struct {
	operatorIdx int
	config      ObscuroConfig
	nodeType    common.NodeType
	l1Data      *params.L1SetupData
	l1Client    ethadapter.EthClient
	logger      gethlog.Logger

	host              *hostcontainer.HostContainer
	enclave           *enclavecontainer.EnclaveContainer
	l1Wallet          wallet.Wallet
	enclaveDBFilepath string
	hostDBFilepath    string
}

func (n *InMemNodeOperator) StopHost() error {
	err := n.host.Stop()
	if err != nil {
		return fmt.Errorf("unable to stop host - %w", err)
	}
	return nil
}

func (n *InMemNodeOperator) Start() error {
	err := n.StartEnclave()
	if err != nil {
		return fmt.Errorf("failed to start enclave - %w", err)
	}
	err = n.StartHost()
	if err != nil {
		return fmt.Errorf("failed to start host - %w", err)
	}
	return nil
}

// StartHost starts the host process in a new thread
func (n *InMemNodeOperator) StartHost() error {
	// even if host was running previously we recreate the container to ensure state is like a new process
	// todo (@matt) - check if host is still running, stop or error?
	n.host = n.createHostContainer()
	go func() {
		err := n.host.Start()
		if err != nil {
			// todo: rework this to return error but able to start all hosts simultaneously
			panic(err)
		}
	}()
	return nil
}

// StartEnclave starts the enclave process in
func (n *InMemNodeOperator) StartEnclave() error {
	// even if enclave was running previously we recreate the container to ensure state is like a new process
	// todo (@matt) - check if enclave is still running?
	n.enclave = n.createEnclaveContainer()
	return n.enclave.Start()
}

func (n *InMemNodeOperator) createHostContainer() *hostcontainer.HostContainer {
	enclavePort := n.config.PortStart + integration.DefaultEnclaveOffset + n.operatorIdx
	enclaveAddr := fmt.Sprintf("%s:%d", network.Localhost, enclavePort)

	p2pPort := n.config.PortStart + integration.DefaultHostP2pOffset + n.operatorIdx
	p2pAddr := fmt.Sprintf("%s:%d", network.Localhost, p2pPort)

	hostConfig := &config.HostConfig{
		ID:                        n.l1Wallet.Address(),
		IsGenesis:                 n.nodeType == common.Sequencer,
		NodeType:                  n.nodeType,
		HasClientRPCHTTP:          true,
		ClientRPCPortHTTP:         uint64(n.config.PortStart + integration.DefaultHostRPCHTTPOffset + n.operatorIdx),
		HasClientRPCWebsockets:    true,
		ClientRPCPortWS:           uint64(n.config.PortStart + integration.DefaultHostRPCWSOffset + n.operatorIdx),
		ClientRPCHost:             network.Localhost,
		EnclaveRPCAddress:         enclaveAddr,
		P2PBindAddress:            p2pAddr,
		P2PPublicAddress:          p2pAddr,
		EnclaveRPCTimeout:         network.EnclaveClientRPCTimeout,
		L1RPCTimeout:              network.DefaultL1RPCTimeout,
		ManagementContractAddress: n.l1Data.MgmtContractAddress,
		MessageBusAddress:         n.l1Data.MessageBusAddr,
		L1ChainID:                 integration.EthereumChainID,
		ObscuroChainID:            integration.TenChainID,
		L1StartHash:               n.l1Data.ObscuroStartBlock,
		SequencerID:               n.config.SequencerID,
		UseInMemoryDB:             false,
		LevelDBPath:               n.hostDBFilepath,
		DebugNamespaceEnabled:     true,
		BatchInterval:             n.config.BatchInterval,
		RollupInterval:            n.config.RollupInterval,
		L1BlockTime:               n.config.L1BlockTime,
		MaxRollupSize:             1024 * 64,
	}

	hostLogger := testlog.Logger().New(log.NodeIDKey, n.l1Wallet.Address(), log.CmpKey, log.HostCmp)

	// create a socket P2P layer
	p2pLogger := hostLogger.New(log.CmpKey, log.P2PCmp)
	svcLocator := host.NewServicesRegistry(n.logger)
	nodeP2p := p2p.NewSocketP2PLayer(hostConfig, svcLocator, p2pLogger, nil)
	// create an enclave client

	enclaveClient := enclaverpc.NewClient(hostConfig, testlog.Logger().New(log.NodeIDKey, n.l1Wallet.Address()))
	rpcServer := clientrpc.NewServer(hostConfig, n.logger)
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&hostConfig.ManagementContractAddress, n.logger)
	l1Repo := l1.NewL1Repository(n.l1Client, []gethcommon.Address{hostConfig.ManagementContractAddress, hostConfig.MessageBusAddress}, n.logger)
	return hostcontainer.NewHostContainer(hostConfig, svcLocator, nodeP2p, n.l1Client, l1Repo, enclaveClient, mgmtContractLib, n.l1Wallet, rpcServer, hostLogger, metrics.New(false, 0, n.logger))
}

func (n *InMemNodeOperator) createEnclaveContainer() *enclavecontainer.EnclaveContainer {
	enclaveLogger := testlog.Logger().New(log.NodeIDKey, n.l1Wallet.Address(), log.CmpKey, log.EnclaveCmp)
	enclavePort := n.config.PortStart + integration.DefaultEnclaveOffset + n.operatorIdx
	enclaveAddr := fmt.Sprintf("%s:%d", network.Localhost, enclavePort)

	hostPort := n.config.PortStart + integration.DefaultHostP2pOffset + n.operatorIdx
	hostAddr := fmt.Sprintf("%s:%d", network.Localhost, hostPort)

	defaultCfg := config.DefaultEnclaveConfig()

	enclaveConfig := &config.EnclaveConfig{
		HostID:                    n.l1Wallet.Address(),
		SequencerID:               n.config.SequencerID,
		HostAddress:               hostAddr,
		Address:                   enclaveAddr,
		NodeType:                  n.nodeType,
		L1ChainID:                 integration.EthereumChainID,
		ObscuroChainID:            integration.TenChainID,
		ValidateL1Blocks:          false,
		WillAttest:                false,
		GenesisJSON:               nil,
		UseInMemoryDB:             false,
		ManagementContractAddress: n.l1Data.MgmtContractAddress,
		MinGasPrice:               big.NewInt(1),
		MessageBusAddress:         n.l1Data.MessageBusAddr,
		SqliteDBPath:              n.enclaveDBFilepath,
		DebugNamespaceEnabled:     true,
		MaxBatchSize:              1024 * 32,
		MaxRollupSize:             1024 * 64,
		BaseFee:                   defaultCfg.BaseFee, // todo @siliev:: fix test transaction builders so this can be different
		GasLimit:                  defaultCfg.GasLimit,
		GasPaymentAddress:         defaultCfg.GasPaymentAddress,
	}
	return enclavecontainer.NewEnclaveContainerWithLogger(enclaveConfig, enclaveLogger)
}

func (n *InMemNodeOperator) Stop() error {
	err := n.host.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop host - %w", err)
	}
	err = n.enclave.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop enclave - %w", err)
	}
	return nil
}

func (n *InMemNodeOperator) HostRPCAddress() string {
	hostPort := n.config.PortStart + integration.DefaultHostRPCWSOffset + n.operatorIdx
	return fmt.Sprintf("ws://%s:%d", network.Localhost, hostPort)
}

func (n *InMemNodeOperator) StopEnclave() error {
	err := n.enclave.Stop()
	if err != nil {
		n.logger.Error("failed to stop enclave - %w", err)

		// try again
		err := n.enclave.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop enclave after second attempt - %w", err)
		}
	}
	return nil
}

func NewInMemNodeOperator(operatorIdx int, config ObscuroConfig, nodeType common.NodeType, l1Data *params.L1SetupData,
	l1Client ethadapter.EthClient, l1Wallet wallet.Wallet, logger gethlog.Logger,
) *InMemNodeOperator {
	// todo (@matt) - put sqlite and levelDB storage in the same temp dir
	sqliteDBPath, err := sqlite.CreateTempDBFile()
	if err != nil {
		panic("failed to create temp sqlite db path")
	}
	levelDBPath, err := os.MkdirTemp("", "levelDB_*")
	if err != nil {
		panic("failed to create temp levelDBPath")
	}

	l1Nonce, err := l1Client.Nonce(l1Wallet.Address())
	if err != nil {
		panic("failed to get l1 nonce")
	}
	l1Wallet.SetNonce(l1Nonce)

	return &InMemNodeOperator{
		operatorIdx:       operatorIdx,
		config:            config,
		nodeType:          nodeType,
		l1Data:            l1Data,
		l1Client:          l1Client,
		l1Wallet:          l1Wallet,
		logger:            logger,
		enclaveDBFilepath: sqliteDBPath,
		hostDBFilepath:    levelDBPath,
	}
}
