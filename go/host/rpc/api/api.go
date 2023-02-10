package api

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/container"

	gethlog "github.com/ethereum/go-ethereum/log"
	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	APIVersion1             = "1.0"
	APINamespaceObscuro     = "obscuro"
	APINamespaceEth         = "eth"
	APINamespaceObscuroScan = "obscuroscan"
	APINamespaceNetwork     = "net"
	APINamespaceTest        = "test"
)

func SequencerAPIs(container container.Container, host commonhost.Host, logger gethlog.Logger) []rpc.API {
	return append(
		baseAPIs(container, host, logger),
		rpc.API{
			Namespace: APINamespaceEth,
			Version:   APIVersion1,
			Service:   NewSequencerEthAPI(host, logger),
			Public:    true,
		},
	)
}

func ValidatorAPIs(container container.Container, host commonhost.Host, logger gethlog.Logger) []rpc.API {
	return append(
		baseAPIs(container, host, logger),
		rpc.API{
			Namespace: APINamespaceEth,
			Version:   APIVersion1,
			Service:   NewValidatorEthAPI(host, logger),
			Public:    true,
		},
	)
}

func baseAPIs(container container.Container, host commonhost.Host, logger gethlog.Logger) []rpc.API {
	return []rpc.API{
		{
			Namespace: APINamespaceObscuro,
			Version:   APIVersion1,
			Service:   NewObscuroAPI(host),
			Public:    true,
		},
		{
			Namespace: APINamespaceObscuroScan,
			Version:   APIVersion1,
			Service:   NewObscuroScanAPI(host),
			Public:    true,
		},
		{
			Namespace: APINamespaceNetwork,
			Version:   APIVersion1,
			Service:   NewNetworkAPI(host),
			Public:    true,
		},
		{
			Namespace: APINamespaceTest,
			Version:   APIVersion1,
			Service:   NewTestAPI(container),
			Public:    true,
		},
		{
			Namespace: APINamespaceEth,
			Version:   APIVersion1,
			Service:   NewFilterAPI(host, logger),
			Public:    true,
		},
	}
}
