package txhandler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type txHandlerImpl struct {
	handlers    []ContractHandler
	mgmtHandler ContractHandler
}

func NewTransactionHandler(mgmtHandler ContractHandler, handlers ...ContractHandler) TxHandler {
	return &txHandlerImpl{
		mgmtHandler: mgmtHandler,
		handlers:    handlers,
	}
}

func (t *txHandlerImpl) PackTx(tx obscurocommon.L1Transaction, from common.Address, nonce uint64) (types.TxData, error) {
	// obscuro will only ever pack transactions using the mgmt contract (for now)
	return t.mgmtHandler.Pack(tx, from, nonce)
}

func (t *txHandlerImpl) UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction {
	// ignore value transfers or contract creations
	if tx.To() == nil {
		log.Trace("UnpackTx: Ignoring transaction %+v", tx)
		return nil
	}

	if tx.To().Hex() == t.mgmtHandler.Address().Hex() {
		return t.mgmtHandler.UnPack(tx)
	}

	for _, handler := range t.handlers {
		if tx.To().Hex() == handler.Address().Hex() {
			return handler.UnPack(tx)
		}
	}

	// ignore contract executions that are not in the registered contract handlers
	log.Trace("UnpackTx: Ignoring transaction %+v", tx)
	return nil
}
