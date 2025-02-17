package container

import (
	"fmt"

	"github.com/ten-protocol/go-ten/tools/faucet/faucet"
	"github.com/ten-protocol/go-ten/tools/faucet/webserver"
)

type FaucetContainer struct {
	faucetServer *faucet.Faucet
	webServer    *webserver.WebServer
}

func NewFaucetContainerFromConfig(cfg *faucet.Config) (*FaucetContainer, error) {
	// we connect to the node via HTTP (config HTTPPort must not be the WSPort for the host)
	nodeAddr := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.HTTPPort)

	f, err := faucet.NewFaucet(nodeAddr, cfg.ChainID.Int64(), cfg.PK[2:])
	if err != nil {
		return nil, err
	}
	bindAddress := fmt.Sprintf(":%d", cfg.ServerPort)
	server := webserver.NewWebServer(f, bindAddress, []byte(cfg.JWTSecret), cfg.DefaultFundAmount)

	return NewFaucetContainer(f, server)
}

func NewFaucetContainer(faucetServer *faucet.Faucet, webServer *webserver.WebServer) (*FaucetContainer, error) {
	return &FaucetContainer{
		faucetServer: faucetServer,
		webServer:    webServer,
	}, nil
}

func (c *FaucetContainer) Start() error {
	return c.webServer.Start()
}

func (c *FaucetContainer) Stop() error {
	return c.webServer.Stop()
}
