package container

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/webserver"

	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/config"
)

type ObscuroScanContainer struct {
	backend   *backend.Backend
	webServer *webserver.WebServer
}

func NewObscuroScanContainer(config *config.Config) (*ObscuroScanContainer, error) {
	client, err := rpc.NewNetworkClient(config.NodeHostAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the obscuro node - %w", err)
	}

	obsClient := obsclient.NewObsClient(client)

	scanBackend := backend.NewBackend(obsClient)
	logger := log.New(log.ObscuroscanCmp, int(gethlog.LvlInfo), config.LogPath)
	webServer := webserver.New(scanBackend, config.ServerAddress, logger)

	return &ObscuroScanContainer{
		backend:   backend.NewBackend(obsClient),
		webServer: webServer,
	}, nil
}

func (c *ObscuroScanContainer) Start() error {
	return c.webServer.Start()
}

func (c *ObscuroScanContainer) Stop() error {
	return c.webServer.Stop()
}
