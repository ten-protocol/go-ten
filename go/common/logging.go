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

// ErrorWithID logs a message at ERROR level with the aggregator's identity prepended.
func ErrorWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Error(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}
