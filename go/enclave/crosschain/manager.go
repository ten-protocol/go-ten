package crosschain

import (
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type CrossChainProcessors struct {
	LocalManager  ObscuroCrossChainManager
	RemoteManager MainNetMessageExtractor
}

func NewCrossChainProcessors(
	l1BusAddress *gethcommon.Address,
	l2BusAddress *gethcommon.Address,
	storage db.Storage, /*key *ecdsa.PrivateKey,*/
	chainId *big.Int,
	logger gethlog.Logger,
) *CrossChainProcessors {
	processors := CrossChainProcessors{}
	processors.LocalManager = NewObscuroMessageBusManager(storage, chainId, logger)
	processors.RemoteManager = NewMainNetExtractor(l1BusAddress, processors.LocalManager.GetBusAddress(), storage, logger)
	return &processors
}
