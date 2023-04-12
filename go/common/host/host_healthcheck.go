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
	P2PStatus       *P2PStatus
	L1BlockProvider *L1BlockProviderStatus
}

// P2PStatus is the representation of the Status of the P2P layer
type P2PStatus struct {
	FailedReceivedMessages int64
	FailedSendMessage      int64
	ReceivedMessages       int64
}

// L1BlockProviderStatus represents the status of the L1 Block Provider
type L1BlockProviderStatus struct {
	Error   string
	Healthy bool
}
