package devnetwork

import (
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter"
)

// L1Network represents the L1Network being used for the devnetwork
// (it could be a local geth docker network, a local in-memory network or even a live public L1)
// todo (@matt) - refactor to use the same NodeOperator approach as the L2?
type L1Network interface {
	Prepare() // ensure L1 connectivity (start nodes if required)
	CleanUp() // shut down  nodes if required, clean up connections
	NumNodes() int
	GetClient(i int) ethadapter.EthClient
	GetBlockTime() time.Duration // expected interval between blocks
}
