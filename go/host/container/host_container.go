package container

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

type HostContainer struct {
	host   commonhost.Host
	logger gethlog.Logger
}

// NewHostContainer wires up the components of the Host service and manages its lifecycle/monitors its status
func NewHostContainer(parsedConfig *config.HostInputConfig) *HostContainer {
	cfg := parsedConfig.ToHostConfig()

	// set the Host ID as the Public Key Address
	ethWallet := wallet.NewInMemoryWalletFromConfig(cfg.PrivateKeyString, cfg.L1ChainID, log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath))
	cfg.ID = ethWallet.Address()

	logger := log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath, log.NodeIDKey, cfg.ID)

	fmt.Printf("Starting host with config: %+v\n", cfg)
	logger.Info(fmt.Sprintf("Starting node with config: %+v", cfg))
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.RollupContractAddress, logger)

	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethadapter.NewEthClient(cfg.L1NodeHost, cfg.L1NodeWebsocketPort, cfg.L1RPCTimeout, cfg.ID, logger)
	if err != nil {
		logger.Crit("could not create Ethereum client.", log.ErrKey, err)
	}

	// update the wallet nonce
	nonce, err := l1Client.Nonce(ethWallet.Address())
	if err != nil {
		logger.Crit("could not retrieve Ethereum account nonce.", log.ErrKey, err)
	}
	ethWallet.SetNonce(nonce)

	// set the Host ID as the Public Key Address
	cfg.ID = ethWallet.Address()

	enclaveClient := enclaverpc.NewClient(cfg, logger)
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	aggP2P := p2p.NewSocketP2PLayer(cfg, p2pLogger)
	h := host.NewHost(cfg, aggP2P, l1Client, enclaveClient, ethWallet, mgmtContractLib, logger)

	// todo: initialise rpc server here and store it on the container, host shouldn't know about rpc server
	return &HostContainer{
		host:   h,
		logger: logger,
	}
}

func (h *HostContainer) Start() error {
	fmt.Println("Starting Obscuro host...")
	h.logger.Info("Starting Obscuro host...")
	h.host.Start()
	return nil
}

func (h *HostContainer) Stop() error {
	h.host.Stop()
	return nil
}
