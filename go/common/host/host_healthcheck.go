package host

type HealthCheck struct {
	OverallHealth bool
	*HealthCheckHost
	*HealthCheckEnclave
}

type HealthCheckEnclave struct {
	EnclaveHealthy bool
}

type HealthCheckHost struct {
	P2PStatus *P2PStatus
}

type P2PStatus struct {
	FailedMessageReads               map[string]int64
	FailedMessageDecodes             map[string]int64
	FailedSendMessagesPeerConnection map[string]int64
	FailedSendMessageWrites          map[string]int64
	ReceivedMessages                 map[string]int64
}
