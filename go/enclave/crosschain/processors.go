package crosschain

import (
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// Object that simply contains the cross chain related structures.
type Processors struct {
	LocalManager  ObscuroCrossChainManager
	RemoteManager BlockMessageExtractor
}

// New - Ensure the managers get bootstrapped properly.
func New(
	l1BusAddress *gethcommon.Address,
	storage db.Storage, /*key *ecdsa.PrivateKey,*/
	chainID *big.Int,
	logger gethlog.Logger,
) *Processors {
	processors := Processors{}
	processors.LocalManager = NewObscuroMessageBusManager(storage, chainID, logger)
	processors.RemoteManager = NewBlockMessageExtractor(l1BusAddress, processors.LocalManager.GetBusAddress(), storage, logger)
	return &processors
}

func (c *Processors) Enabled() bool {
	return c.RemoteManager.Enabled()
}
