package profiler

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"net/http"
	_ "net/http/pprof"
)

const (
	DefaultEnclavePort = 6060
	DefaultHostPort    = 6061
)

type Profiler struct {
	port int
}

func NewProfiler(port int) *Profiler {
	return &Profiler{port: port}
}

func (p *Profiler) Start() error {
	go func() {
		address := fmt.Sprintf("0.0.0.0:%d", p.port)
		log.Info("Profiler started @%s ", address)
		log.Info("%v", http.ListenAndServe(address, nil))
	}()
	return nil
}

func (p *Profiler) Stop() error {
	// todo graceful shutdown
	return nil
}
