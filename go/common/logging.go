package common

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/common/log"
)

// LogWithID logs a message at INFO level with the aggregator's identity prepended.
func LogWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Info(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}

// WarnWithID logs a message at WARN level with the aggregator's identity prepended.
func WarnWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Warn(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}

// TraceWithID logs a message at TRACE level with the aggregator's identity prepended.
func TraceWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Trace(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}

// ErrorWithID logs a message at ERROR level with the aggregator's identity prepended.
func ErrorWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Error(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}

// PanicWithID logs a message at PANIC level with the aggregator's identity prepended.
func PanicWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Panic(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}
