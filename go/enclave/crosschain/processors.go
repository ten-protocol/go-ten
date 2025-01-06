package crosschain

import (
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

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

func (c *Processors) GetL2MessageBusAddress() (gethcommon.Address, common.SystemError) {
	address := c.Local.GetBusAddress()
	if address == nil {
		return gethcommon.Address{}, fmt.Errorf("message bus address not initialised")
	}
	return *address, nil
}
