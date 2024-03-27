package host

// HealthStatus is an interface supported by all Services on the host
type HealthStatus interface {
	OK() bool
	Message() string
}

// HealthCheck is the object returned by the host API with the Health of the Node
type HealthCheck struct {
	OverallHealth bool
	Errors        []string
}

// BasicErrHealthStatus is a simple health status implementation, if the ErrMsg is non-empty then OK() returns false
type BasicErrHealthStatus struct {
	ErrMsg string
}

func (l *BasicErrHealthStatus) OK() bool {
	return l.ErrMsg == ""
}

func (l *BasicErrHealthStatus) Message() string {
	return l.ErrMsg
}

type GroupErrsHealthStatus struct {
	Errors []error
}

func (g *GroupErrsHealthStatus) OK() bool {
	return len(g.Errors) == 0
}

func (g *GroupErrsHealthStatus) Message() string {
	msg := ""
	for i, err := range g.Errors {
		if i > 0 {
			msg += ", "
		}
		msg += err.Error()
	}
	return msg
}
