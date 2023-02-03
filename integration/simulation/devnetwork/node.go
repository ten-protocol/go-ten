package devnetwork

import (
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/go/config"
	enclavecontainer "github.com/obscuronet/go-obscuro/go/enclave/container"
	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
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

	host         *hostcontainer.HostContainer
	enclave      *enclavecontainer.EnclaveContainer
	l1Wallet     wallet.Wallet
	sqliteDBPath string
}

func (n *InMemNodeOperator) StopHost() {
	err := n.host.Stop()
	if err != nil {
		n.logger.Error("unable to stop host - %w", err)
	}
}

func (n *InMemNodeOperator) Start() {
	n.StartEnclave()
	n.StartHost()
}

// StartHost starts the host process in a new thread
func (n *InMemNodeOperator) StartHost() {
	// even if host was running previously we recreate the container to ensure state is like a new process
	// todo: check if host is still running, stop or error?
	n.host = n.createHostContainer()
	go func() {
		err := n.host.Start()
		if err != nil {
			panic(err)
		}
	}()
}

// StartEnclave starts the enclave process in
func (n *InMemNodeOperator) StartEnclave() {
	// even if enclave was running previously we recreate the container to ensure state is like a new process
	// todo: check if enclave is still running?
	n.enclave = n.createEnclaveContainer()
	err := n.enclave.Start()
	if err != nil {
		panic(err)
	}
}

func (n *InMemNodeOperator) createHostContainer() *hostcontainer.HostContainer {
	enclavePort := n.config.PortStart + network.DefaultEnclaveOffset + n.operatorIdx
	enclaveAddr := fmt.Sprintf("%s:%d", network.Localhost, enclavePort)

	p2pPort := n.config.PortStart + network.DefaultHostP2pOffset + n.operatorIdx
	p2pAddr := fmt.Sprintf("%s:%d", network.Localhost, p2pPort)

	hostConfig := &config.HostConfig{
		ID:                     getHostID(n.operatorIdx),
		IsGenesis:              n.nodeType == common.Sequencer,
		NodeType:               n.nodeType,
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      uint64(n.config.PortStart + network.DefaultHostRPCHTTPOffset + n.operatorIdx),
		HasClientRPCWebsockets: true,
		ClientRPCPortWS:        uint64(n.config.PortStart + network.DefaultHostRPCWSOffset + n.operatorIdx),
		ClientRPCHost:          network.Localhost,
		EnclaveRPCAddress:      enclaveAddr,
		P2PBindAddress:         p2pAddr,
		P2PPublicAddress:       p2pAddr,
		EnclaveRPCTimeout:      network.EnclaveClientRPCTimeout,
		L1RPCTimeout:           network.DefaultL1RPCTimeout,
		RollupContractAddress:  n.l1Data.MgmtContractAddress,
		L1ChainID:              integration.EthereumChainID,
		ObscuroChainID:         integration.ObscuroChainID,
		L1StartHash:            n.l1Data.ObscuroStartBlock,
	}

	hostLogger := testlog.Logger().New(log.NodeIDKey, n.operatorIdx, log.CmpKey, log.HostCmp)

	// create a socket P2P layer
	p2pLogger := hostLogger.New(log.CmpKey, log.P2PCmp)
	nodeP2p := p2p.NewSocketP2PLayer(hostConfig, p2pLogger, nil)
	// create an enclave client

	enclaveClient := enclaverpc.NewClient(hostConfig, testlog.Logger().New(log.NodeIDKey, n.operatorIdx))
	rpcServer := clientrpc.NewServer(hostConfig, n.logger)
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&hostConfig.RollupContractAddress, n.logger)
	return hostcontainer.NewHostContainer(hostConfig, nodeP2p, n.l1Client, enclaveClient, mgmtContractLib, n.l1Wallet, rpcServer, hostLogger, metrics.New(false, 0, n.logger))
}

func (n *InMemNodeOperator) createEnclaveContainer() *enclavecontainer.EnclaveContainer {
	enclaveLogger := testlog.Logger().New(log.NodeIDKey, n.operatorIdx, log.CmpKey, log.EnclaveCmp)
	enclavePort := n.config.PortStart + network.DefaultEnclaveOffset + n.operatorIdx
	enclaveAddr := fmt.Sprintf("%s:%d", network.Localhost, enclavePort)

	hostPort := n.config.PortStart + network.DefaultHostP2pOffset + n.operatorIdx
	hostAddr := fmt.Sprintf("%s:%d", network.Localhost, hostPort)

	enclaveConfig := config.EnclaveConfig{
		HostID:                    getHostID(n.operatorIdx),
		SequencerID:               getHostID(0),
		HostAddress:               hostAddr,
		Address:                   enclaveAddr,
		NodeType:                  n.nodeType,
		L1ChainID:                 integration.EthereumChainID,
		ObscuroChainID:            integration.ObscuroChainID,
		ValidateL1Blocks:          false,
		WillAttest:                false,
		GenesisJSON:               nil,
		UseInMemoryDB:             false,
		ManagementContractAddress: n.l1Data.MgmtContractAddress,
		MinGasPrice:               big.NewInt(1),
		MessageBusAddress:         *n.l1Data.MessageBusAddr,
		SqliteDBPath:              n.sqliteDBPath,
	}
	return enclavecontainer.NewEnclaveContainerWithLogger(enclaveConfig, enclaveLogger)
}

func (n *InMemNodeOperator) Stop() {
	err := n.host.Stop()
	if err != nil {
		n.logger.Error("failed to stop host - %w", err)
	}
	err = n.enclave.Stop()
	if err != nil {
		n.logger.Error("failed to stop enclave - %w", err)
	}
}

func (n *InMemNodeOperator) HostRPCAddress() string {
	hostPort := n.config.PortStart + network.DefaultHostRPCWSOffset + n.operatorIdx
	return fmt.Sprintf("ws://%s:%d", network.Localhost, hostPort)
}

func (n *InMemNodeOperator) StopEnclave() {
	err := n.enclave.Stop()
	if err != nil {
		n.logger.Error("failed to stop enclave - %w", err)

		// try again
		err := n.enclave.Stop()
		if err != nil {
			panic("failed to stop enclave after second attempt - " + err.Error())
		}
	}
}

func NewInMemNodeOperator(operatorIdx int, config ObscuroConfig, nodeType common.NodeType, l1Data *params.L1SetupData,
	l1Client ethadapter.EthClient, l1Wallet wallet.Wallet, logger gethlog.Logger,
) *InMemNodeOperator {
	dbFile, err := sql.CreateTempDBFile()
	if err != nil {
		panic("failed to create temp db path")
	}
	return &InMemNodeOperator{
		operatorIdx:  operatorIdx,
		config:       config,
		nodeType:     nodeType,
		l1Data:       l1Data,
		l1Client:     l1Client,
		l1Wallet:     l1Wallet,
		logger:       logger,
		sqliteDBPath: dbFile,
	}
}

func getHostID(nodeIdx int) gethcommon.Address {
	return gethcommon.BigToAddress(big.NewInt(int64(nodeIdx)))
}
