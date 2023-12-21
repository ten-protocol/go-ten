package container

import (
	"fmt"

	"github.com/ten-protocol/go-ten/tools/tenscan/backend"
	"github.com/ten-protocol/go-ten/tools/tenscan/backend/config"
	"github.com/ten-protocol/go-ten/tools/tenscan/backend/webserver"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
)

type TenScanContainer struct {
	backend   *backend.Backend
	webServer *webserver.WebServer
}

func NewTenScanContainer(config *config.Config) (*TenScanContainer, error) {
	client, err := rpc.NewNetworkClient(config.NodeHostAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the obscuro node - %w", err)
	}

	obsClient := obsclient.NewObsClient(client)

	scanBackend := backend.NewBackend(obsClient)
	logger := log.New(log.TenscanCmp, int(gethlog.LvlInfo), config.LogPath)
	webServer := webserver.New(scanBackend, config.ServerAddress, logger)

	logger.Info("Created Obscuro Scan with the following: ", "args", config)
	return &TenScanContainer{
		backend:   backend.NewBackend(obsClient),
		webServer: webServer,
	}, nil
}

func (c *TenScanContainer) Start() error {
	return c.webServer.Start()
}

func (c *TenScanContainer) Stop() error {
	return c.webServer.Stop()
}
