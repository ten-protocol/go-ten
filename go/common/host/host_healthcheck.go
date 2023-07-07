package host

// HealthStatus is an interface supported by all Services on the host
type HealthStatus interface {
	OK() bool
	Message() string
}

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
	L1Repo    bool
	L1Synced  bool
}

// P2PStatus is the representation of the Status of the P2P layer
type P2PStatus struct {
	FailedReceivedMessages int64
	FailedSendMessage      int64
	ReceivedMessages       int64
}

// BasicErrHealthStatus represents the status of a service, if the ErrMsg is non-empty then it reports as "not OK"
type BasicErrHealthStatus struct {
	ErrMsg string
}

func (l *BasicErrHealthStatus) OK() bool {
	return l.ErrMsg == ""
}

func (l *BasicErrHealthStatus) Message() string {
	return l.ErrMsg
}
