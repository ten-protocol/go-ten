package p2p

import (
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	_failedMessageRead        = "msg/inbound/failed_read"
	_failedMessageDecode      = "msg/inbound/failed_decode"
	_failedConnectSendMessage = "msg/outbound/failed_peer_connect"
	_failedWriteSendMessage   = "msg/outbound/failed_write"
	_receivedMessage          = "msg/inbound/success_received"
)

// metrics represents the metrics for the p2p library
type metrics struct {
	hostBasedGauges map[string]*perHostMetrics // gauges broken down per host
}

// newP2PMetrics creates the P2P metrics used by the P2P layer
func newP2PMetrics(registry gethmetrics.Registry) *metrics {
	return &metrics{
		hostBasedGauges: map[string]*perHostMetrics{
			_failedMessageRead:        newperHostMetricMap(registry, _failedMessageRead),
			_failedMessageDecode:      newperHostMetricMap(registry, _failedMessageDecode),
			_failedConnectSendMessage: newperHostMetricMap(registry, _failedConnectSendMessage),
			_failedWriteSendMessage:   newperHostMetricMap(registry, _failedWriteSendMessage),
			_receivedMessage:          newperHostMetricMap(registry, _receivedMessage),
		},
	}
}

// incrementHostEvent adds one (1) to the given instrument at a given host
func (m *metrics) incrementHostEvent(instrument string, host string) {
	m.hostBasedGauges[instrument].inc(host, 1)
}

// status returns the current status of the p2p layer
func (m *metrics) status() *hostcommon.P2PStatus {
	status := &hostcommon.P2PStatus{
		FailedReceivedMessages: int64(0),
		FailedSendMessage:      int64(0),
		ReceivedMessages:       int64(0),
	}

	for gaugeName, gauge := range m.hostBasedGauges {
		switch gaugeName {
		case _receivedMessage:
			status.ReceivedMessages = gauge.totals()
		case _failedMessageRead:
		case _failedMessageDecode:
			status.FailedReceivedMessages += gauge.totals()
		case _failedWriteSendMessage:
		case _failedConnectSendMessage:
			status.FailedSendMessage += gauge.totals()
		}
	}
	return status
}
