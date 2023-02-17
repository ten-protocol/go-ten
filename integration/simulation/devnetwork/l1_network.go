package devnetwork

import (
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

// L1Network represents the L1Network being used for the devnetwork
// (it could be a local geth docker network, a local in-memory network or even a live public L1)
// todo: refactor to use the same NodeOperator approach as the L2?
type L1Network interface {
	Start()
	Stop()
	NumNodes() int
	GetClient(i int) ethadapter.EthClient
	ObscuroSetupData() *params.L1SetupData
}
