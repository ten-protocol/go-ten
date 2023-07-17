package container

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/config"
)

type ObscuroScanContainer struct {
	backend *backend.Backend
}

func NewObscuroScanContainer(config *config.Config) (*ObscuroScanContainer, error) {
	client, err := rpc.NewNetworkClient(config.NodeHostAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the obscuro node - %w", err)
	}

	obsClient := obsclient.NewObsClient(client)

	return &ObscuroScanContainer{
		backend: backend.NewBackend(obsClient),
	}, nil
}

func (c *ObscuroScanContainer) Start() error {
	return nil
}
