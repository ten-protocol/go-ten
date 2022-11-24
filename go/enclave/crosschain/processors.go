package crosschain

import (
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type Processors struct {
	LocalManager  ObscuroCrossChainManager
	RemoteManager MainNetMessageExtractor
}

func New(
	l1BusAddress *gethcommon.Address,
	storage db.Storage, /*key *ecdsa.PrivateKey,*/
	chainID *big.Int,
	logger gethlog.Logger,
) *Processors {
	processors := Processors{}
	processors.LocalManager = NewObscuroMessageBusManager(storage, chainID, logger)
	processors.RemoteManager = NewMainNetExtractor(l1BusAddress, processors.LocalManager.GetBusAddress(), storage, logger)
	return &processors
}

func (c *Processors) Enabled() bool {
	return c.RemoteManager.Enabled()
}
