package host

import (
	gethlog "github.com/ethereum/go-ethereum/log"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

// ServiceFactory is a function that creates a host service
type ServiceFactory[S hostcommon.Service] func(config *config.HostConfig, serviceLocator ServiceLocator, logger gethlog.Logger) (S, error)

type P2PService interface {
	hostcommon.P2P
	hostcommon.Service
}

type L1BlockRepositoryService interface {
	hostcommon.L1BlockRepository
	hostcommon.Service
}

type L1PublisherService interface {
	hostcommon.L1Publisher
	hostcommon.Service
}

type L2BatchRepositoryService interface {
	hostcommon.L2BatchRepository
	hostcommon.Service
}

type EnclaveHostService interface {
	hostcommon.EnclaveService
	hostcommon.Service
}

type LogSubscriptionManagerService interface {
	hostcommon.LogSubscriptionManager
	hostcommon.Service
}

type MetricsService interface {
	hostcommon.Metrics
	hostcommon.Service
}

type RPCServerService interface {
	hostcommon.RPCServer
	hostcommon.Service
}

type P2PLocator interface {
	P2P() hostcommon.P2P
}

type L1RepoLocator interface {
	L1Repo() hostcommon.L1BlockRepository
}

type L1PublisherLocator interface {
	L1Publisher() hostcommon.L1Publisher
}

type L2RepoLocator interface {
	L2Repo() hostcommon.L2BatchRepository
}

type EnclaveLocator interface {
	Enclave() hostcommon.EnclaveService
}

type LogSubsLocator interface {
	LogSubs() hostcommon.LogSubscriptionManager
}

type MetricsLocator interface {
	Metrics() hostcommon.Metrics
}

// DBLocator provides access to the host's database
// Note: we should aim to remove this interface, data should be provided by services and any services that need direct
// DB access should receive a DB connection when they are created
type DBLocator interface {
	DB() *db.DB
}

type HostControlsLocator interface { //nolint:revive
	HostControls() hostcommon.HostControls
}

type ServiceLocator interface {
	P2PLocator
	L1RepoLocator
	L1PublisherLocator
	L2RepoLocator
	EnclaveLocator
	LogSubsLocator
	DBLocator
	MetricsLocator
	HostControlsLocator
}

// DBServiceFactory creates a DB service - this is here because of a circular dependency between the host and the DB
// todo (@matt) change DB to be injected to services that require it rather than a service of its own
func DBServiceFactory(cfg *config.HostConfig, _ ServiceLocator, logger gethlog.Logger) (*db.DB, error) {
	db, err := db.CreateDBFromConfig(cfg, logger)
	if err != nil {
		return nil, err
	}
	return db, nil
}
