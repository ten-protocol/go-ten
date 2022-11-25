package profiler

import (
	"fmt"
	"net/http"

	gethlog "github.com/ethereum/go-ethereum/log"

	_ "net/http/pprof" //nolint:gosec
)

const (
	DefaultEnclavePort = 6060
	DefaultHostPort    = 6061
)

// Profiler stores the data relevant to the profiler instance
type Profiler struct {
	port   int
	logger gethlog.Logger
}

// NewProfiler returns a new profiler that binds on 0.0.0.0:port
func NewProfiler(port int, logger gethlog.Logger) *Profiler {
	return &Profiler{port: port, logger: logger}
}

// Start starts the profiler
func (p *Profiler) Start() error {
	go func() {
		address := fmt.Sprintf("0.0.0.0:%d", p.port)
		p.logger.Info(fmt.Sprintf("Profiler started @%s ", address))
		p.logger.Info(fmt.Sprintf("%v", http.ListenAndServe(address, nil))) //nolint:gosec
	}()
	return nil
}

// Stop stops the profiler
func (p *Profiler) Stop() error {
	// todo graceful shutdown
	return nil
}
