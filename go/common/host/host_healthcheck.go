package host

// HealthCheck is the object returned by the API with the Health and Status of the Node
type HealthCheck struct {
	OverallHealth bool
	*HealthCheckHost
	*HealthCheckEnclave
}

// HealthCheckEnclave is the representation of the Health and Status of the Enclave
type HealthCheckEnclave struct {
	EnclaveHealthy bool
}

// HealthCheckHost is the representation of the Health and Status of the Host
type HealthCheckHost struct {
	P2PStatus *P2PStatus
}

// P2PStatus is the representation of the Status of the P2P layer
type P2PStatus struct {
	FailedMessageReads               map[string]int64
	FailedMessageDecodes             map[string]int64
	FailedSendMessagesPeerConnection map[string]int64
	FailedSendMessageWrites          map[string]int64
	ReceivedMessages                 map[string]int64
}
