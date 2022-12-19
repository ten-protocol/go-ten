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
	port     uint64
	logger   gethlog.Logger
}

func NewHostMetrics(port uint64, logger gethlog.Logger) *Metrics {
	gethmetrics.Enabled = true // TODO hook this under a metrics enabling flag ?
	reg := gethmetrics.DefaultRegistry
	return &Metrics{
		P2PMetrics: NewP2PMetrics(reg),
		registry:   reg,
		port:       port,
		logger:     logger,
	}
}

func (m *Metrics) Start() {
	// starts the metric server
	// TODO hook this under a metrics enabling flag ?
	address := fmt.Sprintf("%s:%d", "127.0.0.1", m.port)
	m.logger.Info("HTTP Metric server started at %s", address)
	exp.Setup(address)
}
