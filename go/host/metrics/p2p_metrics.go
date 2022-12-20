package metrics

import (
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	P2PFailedMessageRead        = "msg/inbound/failed_read"
	P2PFailedMessageDecode      = "msg/inbound/failed_decode"
	P2PFailedConnectSendMessage = "msg/outbound/failed_peer_connect"
	P2PFailedWriteSendMessage   = "msg/outbound/failed_write"
	P2PReceivedMessage          = "msg/inbound/success_received"
)

// P2PMetrics represents the metrics for the p2p library
type P2PMetrics struct {
	hostBasedGauges map[string]*PerStringGaugeMap // gauges broken down per host
}

// NewP2PMetrics creates the P2P metrics used by the P2P layer
func NewP2PMetrics(registry gethmetrics.Registry) *P2PMetrics {
	return &P2PMetrics{
		hostBasedGauges: map[string]*PerStringGaugeMap{
			P2PFailedMessageRead:        NewPerStringGaugeMap(registry, P2PFailedMessageRead),
			P2PFailedMessageDecode:      NewPerStringGaugeMap(registry, P2PFailedMessageDecode),
			P2PFailedConnectSendMessage: NewPerStringGaugeMap(registry, P2PFailedConnectSendMessage),
			P2PFailedWriteSendMessage:   NewPerStringGaugeMap(registry, P2PFailedWriteSendMessage),
			P2PReceivedMessage:          NewPerStringGaugeMap(registry, P2PReceivedMessage),
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
		case P2PReceivedMessage:
			status.ReceivedMessages = gauge.totals()
		case P2PFailedMessageRead:
		case P2PFailedMessageDecode:
			status.FailedReceivedMessages += gauge.totals()
		case P2PFailedWriteSendMessage:
		case P2PFailedConnectSendMessage:
			status.FailedSendMessage += gauge.totals()
		}
	}
	return status
}
