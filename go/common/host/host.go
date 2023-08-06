package host

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Host is the half of the Obscuro node that lives outside the enclave.
type Host interface {
	// Start initializes the main loop of the host.
	Start() error
	Stop()

	// HealthCheck returns the health status of the host + enclave + db
	HealthCheck() (*HealthCheck, error)
}

type BlockStream struct {
	Stream <-chan *types.Block // the channel which will receive the consecutive, canonical blocks
	Stop   func()              // function to permanently stop the stream and clean up any associated processes/resources
}

type BatchMsg struct {
	Batches []*common.ExtBatch // The batches being sent.
	IsLive  bool               // true if these batches are being sent as new, false if in response to a p2p request
}
