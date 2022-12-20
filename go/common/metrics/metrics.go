package metrics

import (
	"fmt"

	"github.com/ethereum/go-ethereum/metrics/exp"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

// Metrics provides the metrics for the host
// it registers the gethmetrics Registry
// and handles the metrics server
type Metrics struct {
	registry gethmetrics.Registry
	port     uint
	logger   gethlog.Logger
}

func New(enabled bool, port uint, logger gethlog.Logger) *Metrics {
	gethmetrics.Enabled = enabled
	return &Metrics{
		registry: gethmetrics.NewRegistry(),
		port:     port,
		logger:   logger,
	}
}

// Start starts the metrics server
func (m *Metrics) Start() {
	// metrics not enabled
	if !gethmetrics.Enabled {
		return
	}
	// starts the metric server
	address := fmt.Sprintf("%s:%d", "0.0.0.0", m.port)
	m.logger.Info("HTTP Metric server started at %s", address)
	// TODO re-write this http server so to have a stop method
	exp.Setup(address)
}

func (m *Metrics) Registry() gethmetrics.Registry {
	return m.registry
}

func (m *Metrics) Stop() {
	// TODO re-write this http server so to have a stop method
}
