package tracers

import (
	"encoding/json"
	
	gethlogger "github.com/ethereum/go-ethereum/eth/tracers/logger"
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*gethlogger.Config
	Tracer  *string
	Timeout *string
	Reexec  *uint64
	// Config specific to given tracer. Note struct logger
	// config are historically embedded in main object.
	TracerConfig json.RawMessage
}
