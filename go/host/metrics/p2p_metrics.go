package metrics

import (
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	FailedMessageRead        = "Failed Message Reads"
	FailedMessageDecode      = "Failed Message Decodes"
	FailedConnectSendMessage = "Failed Peer Connects"
	FailedWriteSendMessage   = "Failed Socket Writes"
	ReceivedMessage          = "Received Messages"
)

// P2PMetrics represents the metrics for the p2p library
type P2PMetrics struct {
	hostBasedGauges map[string]*PerStringGaugeMap // gauges broken down per host
}

// NewP2PMetrics creates the P2P metrics used by the P2P layer
func NewP2PMetrics(registry gethmetrics.Registry) *P2PMetrics {
	return &P2PMetrics{
		hostBasedGauges: map[string]*PerStringGaugeMap{
			FailedMessageRead:        NewPerStringGaugeMap(registry, FailedMessageRead),
			FailedMessageDecode:      NewPerStringGaugeMap(registry, FailedMessageDecode),
			FailedConnectSendMessage: NewPerStringGaugeMap(registry, FailedConnectSendMessage),
			FailedWriteSendMessage:   NewPerStringGaugeMap(registry, FailedWriteSendMessage),
			ReceivedMessage:          NewPerStringGaugeMap(registry, ReceivedMessage),
		},
	}
}

// IncrementHost adds one (1) to the given gauge instrument at a given host
func (m *P2PMetrics) IncrementHost(instrument string, host string) {
	m.hostBasedGauges[instrument].Inc(host, 1)
}

// Status returns the current status of the p2p layer
func (m *P2PMetrics) Status() *hostcommon.P2PStatus {
	status := &hostcommon.P2PStatus{
		FailedReceivedMessages: int64(0),
		FailedSendMessage:      int64(0),
		ReceivedMessages:       int64(0),
	}

	for gaugeName, gauge := range m.hostBasedGauges {
		switch gaugeName {
		case ReceivedMessage:
			status.ReceivedMessages = gauge.totals()
		case FailedMessageRead:
		case FailedMessageDecode:
			status.FailedReceivedMessages += gauge.totals()
		case FailedWriteSendMessage:
		case FailedConnectSendMessage:
			status.FailedSendMessage += gauge.totals()
		}
	}
	return status
}
