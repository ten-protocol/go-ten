package profiler

import (
	"fmt"
	"net/http"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	_ "net/http/pprof" //nolint:gosec
)

const (
	DefaultEnclavePort = 6060
	DefaultHostPort    = 6061
)

type Profiler struct {
	port int
}

// NewProfiler returns a new profiler that will bind on 0.0.0.0:port
func NewProfiler(port int) *Profiler {
	return &Profiler{port: port}
}

// Start starts the profiler
func (p *Profiler) Start() error {
	go func() {
		address := fmt.Sprintf("0.0.0.0:%d", p.port)
		log.Info("Profiler started @%s ", address)
		log.Info("%v", http.ListenAndServe(address, nil))
	}()
	return nil
}

// Stop stops the profiler
func (p *Profiler) Stop() error {
	// todo graceful shutdown
	return nil
}
