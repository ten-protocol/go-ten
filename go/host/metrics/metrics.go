package metrics

import (
	"fmt"

	"github.com/ethereum/go-ethereum/metrics/exp"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

// Metrics provides the metrics for the host
type Metrics struct {
	*P2PMetrics
	registry gethmetrics.Registry
	port     uint
	logger   gethlog.Logger
}

func NewHostMetrics(enabled bool, port uint, logger gethlog.Logger) *Metrics {
	gethmetrics.Enabled = enabled
	reg := gethmetrics.DefaultRegistry
	return &Metrics{
		P2PMetrics: NewP2PMetrics(reg),
		registry:   reg,
		port:       port,
		logger:     logger,
	}
}

func (m *Metrics) Start() {
	// metrics not enabled
	if !gethmetrics.Enabled {
		return
	}
	// starts the metric server
	address := fmt.Sprintf("%s:%d", "0.0.0.0", m.port)
	m.logger.Info("HTTP Metric server started at %s", address)
	exp.Setup(address)
}
