package crosschain

import (
	"math/big"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Processors - contains the cross chain related structures.
type Processors struct {
	Local  Manager
	Remote BlockMessageExtractor
}

func New(
	l1BusAddress *gethcommon.Address,
	storage storage.Storage,
	chainID *big.Int,
	logger gethlog.Logger,
) *Processors {
	processors := Processors{}
	processors.Local = NewObscuroMessageBusManager(storage, chainID, logger)
	processors.Remote = NewBlockMessageExtractor(l1BusAddress, storage, logger)
	return &processors
}

func (c *Processors) Enabled() bool {
	return c.Remote.Enabled()
}
