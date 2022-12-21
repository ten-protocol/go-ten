package metrics

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/metrics/exp"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

var _processCollectionRefreshDuration = 3 * time.Second

// Service provides the metrics for the host
// it registers the gethmetrics Registry
// and handles the metrics server
type Service struct {
	registry gethmetrics.Registry
	port     uint
	logger   gethlog.Logger
}

func New(enabled bool, port uint, logger gethlog.Logger) *Service {
	gethmetrics.Enabled = enabled
	return &Service{
		registry: gethmetrics.NewRegistry(),
		port:     port,
		logger:   logger,
	}
}

// Start starts the metrics server
func (m *Service) Start() {
	// metrics not enabled
	if !gethmetrics.Enabled {
		return
	}

	// TODO make sure to collect process related metrics
	// go gethmetrics.CollectProcessMetrics(_processCollectionRefreshDuration)

	// starts the metric server
	address := fmt.Sprintf("%s:%d", "0.0.0.0", m.port)
	m.logger.Info("HTTP Metric server started at %s", address)
	// TODO re-write this http server so to have a stop method
	exp.Setup(address)
}

// Registry returns the registry for the metrics service
func (m *Service) Registry() gethmetrics.Registry {
	return m.registry
}

func (m *Service) Stop() {
	// TODO re-write this http server so to have a stop method
}
