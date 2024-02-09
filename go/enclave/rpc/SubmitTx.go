package rpc

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

func ExtractSubmitTxRequest(reqParams []any, builder *CallBuilder[common.L2Tx, gethcommon.Hash], _ *EncryptionManager) error {
	l2Tx, err := ExtractTx(reqParams[0].(string))
	if err != nil {
		builder.Err = fmt.Errorf("could not extract transaction. Cause: %w", err)
		return nil
	}
	// we don't return the sender because this call is not authenticated
	builder.Param = l2Tx
	return nil
}

func ExecuteSubmitTx(rpcBuilder *CallBuilder[common.L2Tx, gethcommon.Hash], rpc *EncryptionManager) error {
	if rpc.processors.Local.IsSyntheticTransaction(*rpcBuilder.Param) {
		rpcBuilder.Err = fmt.Errorf("synthetic transaction coming from external rpc")
		return nil
	}

	if err := rpc.service.SubmitTransaction(rpcBuilder.Param); err != nil {
		rpc.logger.Debug("Could not submit transaction", log.TxKey, rpcBuilder.Param.Hash(), log.ErrKey, err)
		rpcBuilder.Err = err
		return nil
	}
	h := rpcBuilder.Param.Hash()
	rpcBuilder.ReturnValue = &h
	return nil
}
