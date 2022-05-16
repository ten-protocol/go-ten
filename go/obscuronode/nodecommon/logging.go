package nodecommon

import (
	"fmt"
	"github.com/obscuronet/obscuro-playground/go/log"
)

// LogWithID logs a message with the aggregator's identity prepended.
func LogWithID(nodeID uint64, msg string, args ...any) {
	formattedMsg := fmt.Sprintf(msg, args)
	log.Log(fmt.Sprintf(">   Agg%d: %s", nodeID, formattedMsg))
}
